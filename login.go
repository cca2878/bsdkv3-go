package bsdkv3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 登录API

const (
	loginAPIPath = "/api/external/login/v3"
	bdInfo       = "cr_nmsl"
)

type loginReq struct {
	baseReqParams
	BdInfo string `form:"bd_info"`
	UserId string `form:"user_id"`
	Pwd    string `form:"pwd"`
}

func (rq loginReq) validate() error {
	return nil
}

type loginResp struct {
	Code json.Number `json:"code"`
	// need_captcha 是否需要验证码。按实际API响应，json字段名need_captch
	NeedCaptcha      optionalValue[string] `json:"need_captch"`
	Nonce            optionalValue[string] `json:"nonce"`
	AccessKey        optionalValue[string] `json:"access_key"`
	Expires          optionalValue[int]    `json:"expires"`
	RealnameVerified optionalValue[int]    `json:"realname_verified"`
	Uid              optionalValue[int]    `json:"uid"`
	Uname            optionalValue[string] `json:"uname"`
	Message          optionalValue[string] `json:"message"`
}

type captLoginReq struct {
	loginReq
	captchaParams
}

func (rq captLoginReq) validate() error {
	return nil
}

type captLoginResp struct {
	loginResp
}

var loginAPI = endpoint[loginReq, loginResp]{
	Method: http.MethodPost,
	Path:   loginAPIPath,
}

var captLoginAPI = endpoint[captLoginReq, captLoginResp]{
	Method: http.MethodPost,
	Path:   loginAPIPath,
}

func (c *Client) hashPwd(pwd string) (string, error) {
	data, err := encryptPKCS1v15(c.publicKey, []byte(c.pwdHash+pwd))
	if err != nil {
		return "", fmt.Errorf("加密密码失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// Login 实现登录功能，支持验证码处理。
func (c *Client) Login(ctx context.Context, u UserInfo) (*SdkAccount, error) {
	pwdHash, err := c.hashPwd(u.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希计算失败: %w", err)
	}

	req := loginReq{
		baseReqParams: c.baseReqParams,
		BdInfo:        bdInfo,
		UserId:        u.Username,
		Pwd:           pwdHash,
	}
	req.Timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)

	resp, err := execAPI(ctx, c, loginAPI, req)
	if err != nil {
		return nil, fmt.Errorf("登录请求错误: %w", err)
	}
	loginResp := resp.Body

	if loginResp.NeedCaptcha.Valid && loginResp.NeedCaptcha.Value == "1" {
		captchaParams, err := c.handleCaptcha(ctx)
		if err != nil {
			return nil, fmt.Errorf("验证码处理失败: %w", err)
		}

		captReq := captLoginReq{
			loginReq:      req,
			captchaParams: *captchaParams,
		}

		captResp, err := execAPI(ctx, c, captLoginAPI, captReq)
		if err != nil {
			return nil, fmt.Errorf("验证码登录请求错误: %w", err)
		}
		if captResp.Body.Code.String() != "0" {
			return nil, fmt.Errorf("验证码登录错误: (%s) %s", captResp.Body.Code.String(), optionalString(captResp.Body.Message))
		}

		if account, ok := sdkAccountFromLoginResp(captResp.Body.loginResp, u); ok {
			return account, nil
		}
		return nil, fmt.Errorf("登录错误: (%s) %s", captResp.Body.Code.String(), optionalString(captResp.Body.Message))
	}

	if account, ok := sdkAccountFromLoginResp(loginResp, u); ok {
		return account, nil
	}
	return nil, fmt.Errorf("登录错误: (%s) %s", loginResp.Code.String(), optionalString(loginResp.Message))
}

func optionalString(o optionalValue[string]) string {
	if o.Valid {
		return o.Value
	}
	return ""
}

func sdkAccountFromLoginResp(r loginResp, u UserInfo) (*SdkAccount, bool) {
	if r.Code.String() != "0" || !r.AccessKey.Valid || r.AccessKey.Value == "" {
		return nil, false
	}
	uid := ""
	if r.Uid.Valid {
		uid = strconv.Itoa(r.Uid.Value)
	}
	return &SdkAccount{
		AccessKey: r.AccessKey.Value,
		Uid:       uid,
		Platform:  u.Platform,
		Channel:   u.Channel,
	}, true
}
