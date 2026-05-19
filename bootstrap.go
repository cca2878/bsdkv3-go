package bsdkv3

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 取外部配置

const extConfAPIPath = "/api/external/config/v3"

type extConfReq struct {
	baseReqParams
}

func (rq extConfReq) validate() error {
	return nil
}

type extConfResp struct {
	ConfigLoginHttps   optionalValue[string] `json:"config_login_https"`
	ConfigAndroidHttps optionalValue[string] `json:"config_login_android_https"`
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
	baseReqParams
	CipherType string `form:"cipher_type"`
}

func (rq getCipherV3Req) validate() error {
	if rq.CipherType == "" {
		return fmt.Errorf("cipher_type不能为空")
	}
	return nil
}

type getCipherV3Resp struct {
	CipherKey optionalValue[string] `json:"cipher_key"`
	Hash      optionalValue[string] `json:"hash"`
}

var getCipherV3API = endpoint[getCipherV3Req, getCipherV3Resp]{
	Method: http.MethodPost,
	Path:   getCipherV3APIPath,
}

// bootstrap 拉取外部配置与登录密钥，完成 Client 初始化。
func (c *Client) bootstrap(ctx context.Context) error {
	reqConf := extConfReq{
		baseReqParams: c.baseReqParams,
	}
	reqConf.Timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)

	confResp, err := execAPI(ctx, c, extConfAPI, reqConf)
	if err != nil {
		return fmt.Errorf("获取外部配置失败: %w", err)
	}
	if !confResp.Body.ConfigLoginHttps.Valid || confResp.Body.ConfigLoginHttps.Value == "" {
		return fmt.Errorf("获取的登录Hosts配置为空")
	}

	loginHosts := confResp.Body.ConfigLoginHttps.Value
	c.logger.Info("更新登录Hosts")
	c.logger.Debug("%s", loginHosts)
	c.hostMgr.updateHosts(strings.Split(loginHosts, ","))
	// c.hosts.UpdateHosts(config.ParseHostsStr(config.HostTypeLoginHttps, confResp.Body.ConfigLoginHttps))

	cipherResp, err := execAPI(ctx, c, getCipherV3API, getCipherV3Req{
		baseReqParams: c.baseReqParams,
		CipherType:    cipherType,
	})
	if err != nil {
		return fmt.Errorf("获取登录Secrets失败: %w", err)
	}

	if !cipherResp.Body.Hash.Valid || !cipherResp.Body.CipherKey.Valid {
		return fmt.Errorf("获取的登录Secrets响应缺少必要字段")
	}
	c.pwdHash = cipherResp.Body.Hash.Value
	publicKey, err := parsePubkeyFromPEM(cipherResp.Body.CipherKey.Value)
	if err != nil {
		return fmt.Errorf("解析公钥失败: %w", err)
	}
	c.publicKey = *publicKey

	return nil
}
