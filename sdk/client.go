// Package sdk provides a Go client for authentication and validation services.
//
// This package offers a clean, idiomatic Go interface for handling authentication
// workflows, including automatic captcha validation and secure credential handling.
//
// # Quick Start
//
// Create a client and authenticate:
//
//	client, err := sdk.NewBSdkV3Client(config.AppkeyPcr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	user := sdk.UserInfo{
//		Username: "your_username",
//		Password: "your_password",
//	}
//
//	accessKey, err := client.Login(user)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Configuration
//
// The client supports flexible configuration via the config package:
//
//	cfg := config.NewConfig(
//		config.WithRequestTimeout(10 * time.Second),
//		config.WithAppId("custom_app_id"),
//	)
//	client, err := sdk.NewBSdkV3Client(appKey, sdk.WithConfig(cfg))
//
// # Context and Cancellation
//
// The client automatically manages contexts internally. To cancel all pending
// operations, call Close() on the client:
//
//	client.Close() // Cancels all pending requests
package sdk

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
	"strings"

	"github.com/cca2878/bsdkv3-go/sdk/config"
	"github.com/cca2878/bsdkv3-go/sdk/log"
)

// BSdkV3Client 封装 HTTP 客户端和 API 调用。
//
// 客户端提供了认证和验证服务的完整功能，包括：
// - 自动获取和管理配置
// - 密码加密和签名计算
// - 验证码处理
// - 请求签名和发送
//
// 使用示例：
//
//	client, err := NewBSdkV3Client("your_app_key")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	user := UserInfo{Username: "user", Password: "pass"}
//	accessKey, err := client.Login(user)
type BSdkV3Client struct {
	client      *resty.Client
	AccessKey   string // 登录后保存的 token
	formEncoder *form.Encoder
	publicKey   *rsa.PublicKey
	pwdHash     string
	appKey      string
	ctx         context.Context    // 内部创建的 context
	ctxCancel   context.CancelFunc // 用于取消 context
	config      *config.Config     // 客户端配置
}

// ClientOption 定义客户端选项函数类型，用于配置客户端行为
type ClientOption func(*BSdkV3Client)

// WithConfig 设置客户端配置
//
// 使用示例：
//
//	cfg := config.NewConfig(config.WithRequestTimeout(10 * time.Second))
//	client, err := NewBSdkV3Client(appKey, WithConfig(cfg))
func WithConfig(cfg *config.Config) ClientOption {
	return func(c *BSdkV3Client) {
		c.config = cfg
	}
}

// NewBSdkV3Client 创建一个新的 BSdkV3Client 实例。
//
// 该函数会自动：
// - 创建内部 context 用于请求管理
// - 初始化 HTTP 客户端
// - 获取服务器配置和加密密钥
//
// 参数：
//   - appKey: 应用密钥，用于签名计算
//   - options: 可选的配置选项
//
// 返回错误类型：
//   - *ConfigError: 配置获取或初始化失败
//
// 使用示例：
//
//	client, err := NewBSdkV3Client(config.AppkeyPcr)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
func NewBSdkV3Client(appKey string, options ...ClientOption) (*BSdkV3Client, error) {
	// 创建一个带有取消功能的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 使用默认配置创建客户端
	client := &BSdkV3Client{
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
		// 设置默认Content-Type为form-urlencoded
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("User-Agent", "Mozilla/5.0 BSGameSDK").
		SetHeader("cversion", "1").
		// 设置默认超时等
		SetTimeout(client.config.RequestTimeout)

	err := client.getConfig()
	if err != nil {
		return nil, NewConfigError("initialization", err)
	}
	return client, nil
}

// GetConfig 返回客户端配置
func (c *BSdkV3Client) GetConfig() *config.Config {
	return c.config
}

func (c *BSdkV3Client) getConfig() error {
	confReq := NewBSdkV3ExtConfReq(c.config)
	var confResp BSdkV3ExtConfResp
	_, err := c.execReq(c.ctx, confReq, &confResp)
	if err != nil {
		return fmt.Errorf("获取外部配置失败: %w", err)
	}
	if confResp.ConfigLoginHttps == "" {
		return fmt.Errorf("获取的登录HTTPS配置为空")
	}

	log.Info("更新登录Hosts")
	log.Debug("配置信息: %s", confResp.ConfigLoginHttps)
	config.GetHostConfig().UpdateHosts(config.ParseHostsStr(config.HostTypeLoginHttps, confResp.ConfigLoginHttps))

	cipherReq := NewBSdkGetCipherV3Req(c.config)
	var cipherResp BSdkGetCipherV3Resp
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

// calculateSign 计算请求签名
func (c *BSdkV3Client) calculateSign(requestBody interface{}) (string, error) {
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
func (c *BSdkV3Client) hashPwd(pwd string) (string, error) {
	data, err := encryptPKCS1v15(c.publicKey, []byte(c.pwdHash+pwd))
	if err != nil {
		return "", fmt.Errorf("加密密码失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// prepareRequest 准备请求，处理签名和表单数据
func (c *BSdkV3Client) prepareRequest(ctx context.Context, requestBody interface{}) (*resty.Request, error) {
	// 计算签名
	sign, err := c.calculateSign(requestBody)
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
func (c *BSdkV3Client) execReq(ctx context.Context, request Request, result interface{}) (*resty.Response, error) {
	req, err := c.prepareRequest(ctx, request)
	if err != nil {
		log.Error("准备请求失败: %v", err)
		return nil, err
	}

	url, err := request.GetUrl()
	if err != nil {
		log.Error("获取URL失败: %v", err)
		return nil, err
	}

	log.Debug("发送请求: %s", url.String())

	// 发送请求并处理结果
	resp, err := req.SetResult(result).Post(url.String())
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

// Login 执行用户登录流程，支持自动验证码处理。
//
// 该方法会自动处理以下步骤：
// 1. 密码加密
// 2. 发送登录请求
// 3. 检查是否需要验证码
// 4. 如需要，自动处理验证码验证
// 5. 保存访问令牌
//
// 参数：
//   - u: 用户登录信息，包含用户名和密码
//
// 返回：
//   - *string: 成功时返回访问令牌
//   - error: 失败时返回具体错误，可能的类型有 *AuthError
//
// 使用示例：
//
//	user := UserInfo{Username: "user", Password: "pass"}
//	accessKey, err := client.Login(user)
//	if err != nil {
//		var authErr *AuthError
//		if errors.As(err, &authErr) {
//			log.Printf("认证失败 [%s]: %s", authErr.Code, authErr.Message)
//		}
//		return err
//	}
func (c *BSdkV3Client) Login(u UserInfo) (*string, error) {
	// 使用内部创建的 context
	ctx := c.ctx
	// Hash密码
	pwdHash, err := c.hashPwd(u.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希计算失败: %w", err)
	}

	user := UserInfo{
		Username: u.Username,
		Password: pwdHash,
	}

	// 构造登录请求
	loginReq := NewBSdkV3LoginReq(c.config, user)

	// 发起第一次登录请求
	var loginResp BSdkV3LoginResp
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
		captLoginReq := NewBSdkV3CaptLoginReq(c.config, user, *captchaParams)

		// 发起带验证码的登录请求
		var captLoginResp BSdkV3CaptLoginResp

		_, err = c.execReq(ctx, captLoginReq, &captLoginResp)

		if err != nil {
			return nil, fmt.Errorf("验证码登录请求错误: %w", err)
		}
		if captLoginResp.Code.String() != "0" {
			return nil, NewAuthError(captLoginResp.Code.String(), *captLoginResp.Message, nil)
		}

		// 保存token（如果有的话）
		if captLoginResp.AccessKey != nil {
			c.AccessKey = *captLoginResp.AccessKey
			c.client.SetAuthToken(c.AccessKey)
		}
	} else
	// 不需要验证码，直接使用第一次登录的结果
	if loginResp.Code.String() == "0" && loginResp.AccessKey != nil {
		c.AccessKey = *loginResp.AccessKey
		c.client.SetAuthToken(c.AccessKey)
	} else {
		return nil, NewAuthError(loginResp.Code.String(), *loginResp.Message, nil)
	}

	return &c.AccessKey, nil
}

func (c *BSdkV3Client) handleCaptcha(ctx context.Context) (*CaptchaParams, error) {
	// 请求验证码参数
	captchaReq := NewBSdkStartCaptchaReq(c.config)

	// 准备请求，处理签名和表单数据
	req, err := c.prepareRequest(ctx, captchaReq)
	if err != nil {
		return nil, fmt.Errorf("准备验证码请求失败: %w", err)
	}

	url, err := captchaReq.GetUrl()
	if err != nil {
		return nil, fmt.Errorf("获取验证码URL失败: %w", err)
	}

	var captchaResp BSdkStartCaptchaResp
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
	captchaParams := CaptchaParams{
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
func (c *BSdkV3Client) startCaptcha(ctx context.Context) (*BSdkStartCaptchaResp, error) {
	// 构造请求体
	captchaReq := NewBSdkStartCaptchaReq(c.config)

	var result BSdkStartCaptchaResp
	// 发起 POST 请求
	_, err := c.execReq(ctx, captchaReq, &result)

	if err != nil {
		return nil, fmt.Errorf("获取验证码请求错误: %w", err)
	}

	return &result, nil
}

// Close 关闭客户端并取消所有未完成的请求
func (c *BSdkV3Client) Close() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}
}
