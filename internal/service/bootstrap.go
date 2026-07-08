package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/cca2878/bsdkv3-go/internal/base"
	"github.com/cca2878/bsdkv3-go/internal/interceptor"
)

// 取外部配置

const (
	defaultInitConfHost = "https://p.biligame.com"
	extConfAPIPath      = "/api/external/config/v3"
)

type extConfReq struct {
}

func (rq extConfReq) validate() error {
	return nil
}

type extConfResp struct {
	ConfigLoginHttps   base.OptionalValue[string] `json:"config_login_https"`
	ConfigAndroidHttps base.OptionalValue[string] `json:"config_login_android_https"`
}

var extConfAPI = endpoint[extConfReq, extConfResp]{
	Method: http.MethodPost,
	Path:   defaultInitConfHost + extConfAPIPath,
}

// 登录Secrets API

const (
	getCipherV3APIPath = "/api/external/issue/cipher/v3"
	cipherType         = "bili_login_rsa"
)

type getCipherV3Req struct {
	CipherType string `form:"cipher_type"`
}

func (rq getCipherV3Req) validate() error {
	if rq.CipherType == "" {
		return fmt.Errorf("cipher_type不能为空")
	}
	return nil
}

type getCipherV3Resp struct {
	CipherKey base.OptionalValue[string] `json:"cipher_key"`
	Hash      base.OptionalValue[string] `json:"hash"`
}

var getCipherV3API = endpoint[getCipherV3Req, getCipherV3Resp]{
	Method: http.MethodPost,
	Path:   getCipherV3APIPath,
}

type CipherResult struct {
	PublicKey rsa.PublicKey
	HashSalt  string
}

// FetchBootstrapHosts 专门负责去服务器拉取初始 Hosts 列表
func FetchBootstrapHosts(ctx context.Context, invoker interceptor.Invoker) ([]string, error) {
	confResp, err := execAPI(ctx, invoker, extConfAPI, extConfReq{})
	if err != nil {
		return nil, fmt.Errorf("获取外部配置失败: %w", err)
	}

	if !confResp.ConfigLoginHttps.Valid || confResp.ConfigLoginHttps.Value == "" {
		return nil, fmt.Errorf("获取的登录Hosts配置为空")
	}

	loginHosts := confResp.ConfigLoginHttps.Value
	var hosts []string
	for h := range strings.SplitSeq(loginHosts, ",") {
		h = strings.TrimSpace(h)
		if h != "" {
			hosts = append(hosts, h)
		}
	}
	if len(hosts) == 0 {
		return nil, fmt.Errorf("解析的登录Hosts列表为空")
	}

	return hosts, nil
}

// FetchCipher 专门负责去服务器拉取登录 Secrets
// 注意：此时传入的 invoker 应该是包装了高可用 HostMgr 和 RetryInterceptor 的，
// 因此此处直接发出相对路径请求即可享受高可用、重试和路由惩罚机制。
func FetchCipher(ctx context.Context, invoker interceptor.Invoker) (*CipherResult, error) {
	cipherResp, err := execAPI(ctx, invoker, getCipherV3API, getCipherV3Req{
		CipherType: cipherType,
	})
	if err != nil {
		return nil, fmt.Errorf("获取登录Secrets失败: %w", err)
	}

	if !cipherResp.Hash.Valid || !cipherResp.CipherKey.Valid {
		return nil, fmt.Errorf("获取的登录Secrets响应缺少必要字段")
	}
	publicKey, err := parsePubkeyFromPEM(cipherResp.CipherKey.Value)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	return &CipherResult{
		PublicKey: *publicKey,
		HashSalt:  cipherResp.Hash.Value,
	}, nil
}
