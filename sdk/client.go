package bsdkv3

import (
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"bsdkv3-go/sdk/config"
	"bsdkv3-go/sdk/log"
)

// Client 封装 HTTP 客户端和 API 调用
type Client struct {
	client      *resty.Client
	formEncoder *form.Encoder
	publicKey   *rsa.PublicKey
	pwdHash     string
	appKey      string
	ctx         context.Context    // 内部创建的 context
	ctxCancel   context.CancelFunc // 用于取消 context
	config      *config.Config     // 客户端配置
}

// ClientOption 定义客户端选项
type ClientOption func(*Client)

// withConfig 设置客户端配置
func withConfig(cfg *config.Config) ClientOption {
	return func(c *Client) {
		c.config = cfg
	}
}

// NewClient 创建一个新的 Client 实例。
// 内部自动创建 context，调用方无需显式传递 context
func NewClient(appKey string, options ...ClientOption) (*Client, error) {
	// 创建一个带有取消功能的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 使用默认配置创建客户端
	client := &Client{
		formEncoder: form.NewEncoder(),
		appKey:      appKey,
		ctx:         ctx,
		ctxCancel:   cancel,
		config:      config.NewDefaultConfig(), // 默认配置
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	// 创建并配置 HTTP 客户端
	client.client = resty.New().
		// 设置Headers
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Mozilla/5.0 BSGameSDK",
			"cversion":     "1",
		}).
		// debug
		//SetProxy("http://127.0.0.1:8080").
		// 设置默认超时等
		SetTimeout(client.config.RequestTimeout)

	err := client.getConfig()
	if err != nil {
		return nil, fmt.Errorf("配置失败: %w", err)
	}
	return client, nil
}

func (c *Client) getConfig() error {
	confReq := newExtConfReq(c.config)
	var confResp extConfResp
	_, err := c.execReq(c.ctx, confReq, &confResp)
	if err != nil {
		return fmt.Errorf("获取外部配置失败: %w", err)
	}
	if confResp.ConfigLoginHttps == "" {
		return fmt.Errorf("获取的登录HTTPS配置为空")
	}

	log.Info("更新登录Hosts")
	log.Debug(confResp.ConfigLoginHttps)
	config.GetHostConfig().UpdateHosts(config.ParseHostsStr(config.HostTypeLoginHttps, confResp.ConfigLoginHttps))

	cipherReq := newGetCipherV3Req(c.config)
	var cipherResp getCipherV3Resp
	_, err = c.execReq(c.ctx, cipherReq, &cipherResp)
	if err != nil {
		return fmt.Errorf("获取密钥失败: %w", err)
	}

	c.pwdHash = cipherResp.Hash
	c.publicKey, err = parsePublicKeyFromPEM(cipherResp.CipherKey)
	if err != nil {
		return fmt.Errorf("解析公钥失败: %w", err)
	}

	return nil
}

// calcSign 计算请求签名
func (c *Client) calcSign(requestBody interface{}) (string, error) {
	// 请求体编码为map
	values, err := c.formEncoder.Encode(requestBody)
	if err != nil {
		return "", fmt.Errorf("表单编码失败: %w", err)
	}

	// 获取所有键并排序
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按排序后的键顺序拼接值
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(values.Get(k))
	}

	// 拼接AppKey
	sb.WriteString(c.appKey)

	// 计算MD5值作为sign
	sum := md5.Sum([]byte(sb.String()))
	return hex.EncodeToString(sum[:]), nil
}

// hashPwd 对密码进行加密
func (c *Client) hashPwd(pwd string) (string, error) {
	data, err := encryptPKCS1v15(c.publicKey, []byte(c.pwdHash+pwd))
	if err != nil {
		return "", fmt.Errorf("加密密码失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// preReq 准备请求，处理签名和表单数据
func (c *Client) preReq(ctx context.Context, requestBody interface{}) (*resty.Request, error) {
	// 计算签名
	sign, err := c.calcSign(requestBody)
	if err != nil {
		return nil, fmt.Errorf("计算签名失败: %w", err)
	}

	// 使用form库将请求体编码为表单值
	values, err := c.formEncoder.Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("请求体编码失败: %w", err)
	}

	// 添加签名
	values.Set("sign", sign)

	// 创建请求
	req := c.client.R().
		SetContext(ctx).
		SetFormDataFromValues(values)

	return req, nil
}

// execReq
func (c *Client) execReq(ctx context.Context, request iRequest, result any) (*resty.Response, error) {
	req, err := c.preReq(ctx, request)
	if err != nil {
		log.Error("准备请求失败: %v", err)
		return nil, err
	}

	url, err := request.getUrl()
	if err != nil {
		log.Error("获取URL失败: %v", err)
		return nil, err
	}

	log.Debug("发送请求: %s", url.String())

	// 发送请求并处理结果
	resp, err := req.SetResult(result).Execute(request.getMethod(), url.String())
	if err != nil {
		log.Error("请求发送失败: %v", err)
		return resp, err
	}

	log.Debug("收到响应: 状态码=%d, 内容长度=%d", resp.StatusCode(), len(resp.Body()))

	if resp.StatusCode() != http.StatusOK {
		log.Error("请求失败，状态码: %d", resp.StatusCode())
		return resp, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode())
	}

	return resp, nil
}

// Login 方法实现登录功能，支持验证码处理
// 使用内部创建的 context，调用方无需显式传递
func (c *Client) Login(u UserInfo) (*SdkAccount, error) {
	// 使用内部创建的 context
	ctx := c.ctx
	// Hash密码
	pwdHash, err := c.hashPwd(u.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希计算失败: %w", err)
	}

	// 构造登录请求
	loginReq := newLoginReq(c.config).(*loginReq)
	loginReq.UserId = u.Username
	loginReq.Pwd = pwdHash

	// 发起第一次登录请求
	var loginResp loginResp
	_, err = c.execReq(ctx, loginReq, &loginResp)

	if err != nil {
		return nil, fmt.Errorf("登录请求错误: %w", err)
	}

	// 检查是否需要验证码
	if loginResp.NeedCaptcha != nil && *loginResp.NeedCaptcha == "1" {
		// 需要验证码，启动验证流程
		captchaParams, err := c.handleCaptcha(ctx)
		if err != nil {
			return nil, fmt.Errorf("验证码处理失败: %w", err)
		}

		// 构造带验证码的登录请求
		captLoginReq := newCaptLoginReq(c.config, *captchaParams).(*captLoginReq)
		captLoginReq.UserId = u.Username
		captLoginReq.Pwd = pwdHash

		// 发起带验证码的登录请求
		var captLoginResp captLoginResp

		_, err = c.execReq(ctx, captLoginReq, &captLoginResp)

		if err != nil {
			return nil, fmt.Errorf("验证码登录请求错误: %w", err)
		}
		if captLoginResp.Code.String() != "0" {
			return nil, fmt.Errorf("验证码登录错误: (%s) %s", captLoginResp.Code.String(), *captLoginResp.Message)
		}

		// AccessKey（如果有的话）
		if captLoginResp.AccessKey != nil {
			return &SdkAccount{
				AccessKey: *captLoginResp.AccessKey,
				Uid:       strconv.Itoa(*captLoginResp.Uid),
				Platform:  u.Platform,
				Channel:   u.Channel,
			}, nil
		} else {
			return nil, fmt.Errorf("登录错误: (%s) %s", loginResp.Code.String(), *loginResp.Message)
		}
	} else
	// 不需要验证码，直接使用第一次登录的结果
	if loginResp.Code.String() == "0" && loginResp.AccessKey != nil {
		return &SdkAccount{
			AccessKey: *loginResp.AccessKey,
			Uid:       strconv.Itoa(*loginResp.Uid),
			Platform:  u.Platform,
			Channel:   u.Channel,
		}, nil
	} else {
		return nil, fmt.Errorf("登录错误: (%s) %s", loginResp.Code.String(), *loginResp.Message)
	}

}

func (c *Client) handleCaptcha(ctx context.Context) (*captchaParams, error) {
	// 请求验证码参数
	captchaReq := newStartCaptchaReq(c.config)

	// 准备请求，处理签名和表单数据
	req, err := c.preReq(ctx, captchaReq)
	if err != nil {
		return nil, fmt.Errorf("准备验证码请求失败: %w", err)
	}

	url, err := captchaReq.getUrl()
	if err != nil {
		return nil, fmt.Errorf("获取验证码URL失败: %w", err)
	}

	var captchaResp startCaptchaResp
	resp, err := req.
		SetResult(&captchaResp).
		Post(url.String())

	if err != nil {
		return nil, fmt.Errorf("发送验证码请求失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("验证码请求返回非成功状态码: %d", resp.StatusCode())
	}

	// 验证码验证逻辑，暂时仅使用远程验证器
	ret, err := NewRemoteValidator().Validate()
	if err != nil {
		return nil, fmt.Errorf("远程验证码校验失败: %w", err)
	}

	// 构造验证码参数
	captchaParams := captchaParams{
		CaptchaType: "1",
		Validate:    ret.Validate,
		Challenge:   ret.Challenge,
		GtUserId:    ret.GtUserId,
		SecCode:     ret.Validate + "|jordan",
		CToken:      ret.GtUserId,
	}

	return &captchaParams, nil
}

// startCaptcha 方法获取验证码信息
func (c *Client) startCaptcha(ctx context.Context) (*startCaptchaResp, error) {
	// 构造请求体
	captchaReq := newStartCaptchaReq(c.config)

	var result startCaptchaResp
	// 发起 POST 请求
	_, err := c.execReq(ctx, captchaReq, &result)

	if err != nil {
		return nil, fmt.Errorf("获取验证码请求错误: %w", err)
	}

	return &result, nil
}

// Close 关闭客户端并取消所有未完成的请求
func (c *Client) Close() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}
}
