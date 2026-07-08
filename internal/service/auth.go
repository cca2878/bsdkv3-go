package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cca2878/bsdkv3-go/internal/base"
)

// 登录API

const (
	loginAPIPath = "/api/external/login/v3"
	bdInfo       = "cr_nmsl"
)

type loginReq struct {
	BdInfo string `form:"bd_info"`
	UserId string `form:"user_id"`
	Pwd    string `form:"pwd"`

	// 验证码相关参数（如果需要）
	captchaParams
}

type captchaParams struct {
	// 验证码相关参数（如果需要）
	CaptchaType string `form:"captcha_type,omitempty"`
	SecCode     string `form:"seccode,omitempty"`
	Validate    string `form:"validate,omitempty"`
	GtUserId    string `form:"gt_user_id,omitempty"`
	CToken      string `form:"ctoken,omitempty"`
	Challenge   string `form:"challenge,omitempty"`
}

func (rq loginReq) validate() error {
	return nil
}

type loginResp struct {
	Code base.OptionalValue[json.Number] `json:"code"`
	// need_captcha 是否需要验证码。按实际API响应，json字段名need_captch
	NeedCaptcha      base.OptionalValue[string] `json:"need_captch"`
	Nonce            base.OptionalValue[string] `json:"nonce"`
	AccessKey        base.OptionalValue[string] `json:"access_key"`
	Expires          base.OptionalValue[int]    `json:"expires"`
	RealnameVerified base.OptionalValue[int]    `json:"realname_verified"`
	Uid              base.OptionalValue[int]    `json:"uid"`
	Uname            base.OptionalValue[string] `json:"uname"`
	Message          base.OptionalValue[string] `json:"message"`
}

var loginAPI = endpoint[loginReq, loginResp]{
	Method: http.MethodPost,
	Path:   loginAPIPath,
}

func (s *Service) hashPwd(pwd string) (string, error) {
	data, err := encryptPKCS1v15(s.cipher.PublicKey, []byte(s.cipher.HashSalt+pwd))
	if err != nil {
		return "", fmt.Errorf("加密密码失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// UserInfo 登录用户凭据与渠道信息。
type UserInfo struct {
	Username string
	Password string

	Platform string
	Channel  string
}

// SdkAccount 登录成功后返回的账号信息。
type SdkAccount struct {
	Uid       string
	AccessKey string

	Platform string
	Channel  string
}

func (s *Service) Login(ctx context.Context, user UserInfo) (*SdkAccount, error) {
	pwdHash, err := s.hashPwd(user.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希计算失败: %w", err)
	}
	// 初始状态下，没有验证码票据
	var ticket *captchaResult

	// 防御性编程：最多允许触发 2 次验证码挑战，防止服务端死循环下发验证码
	maxAttempts := 2

	for range maxAttempts {
		reqBody := loginReq{
			BdInfo: bdInfo,
			UserId: user.Username,
			Pwd:    pwdHash,
		}

		// 如果上一轮循环拿到了票据，这次就带上它！
		if ticket != nil {
			reqBody.captchaParams = captchaParams(*ticket)
		}

		// 2. 发射请求，穿透 Pipeline！
		respBody, err := execAPI(ctx, s.doer, loginAPI, reqBody)
		if err != nil {
			return nil, err // 物理网络错误，直接抛出
		}

		// 3. 业务状态机判断 (The State Machine)

		// 状态 A：服务器要求进行人机验证
		if respBody.NeedCaptcha.Valid && respBody.NeedCaptcha.Value == "1" {
			captchaParams, err := s.handleCaptcha(ctx)
			if err != nil {
				return nil, fmt.Errorf("人机验证失败或主动取消: %w", err)
			}
			// 验证成功！更新 ticket，然后 continue 进入下一轮发包循环
			ticket = captchaParams
			continue
		}

		// 状态 B：登录失败（密码错误、账号封禁等非验证码错误）
		if respBody.Code.Valid && respBody.Code.Value.String() != "0" {
			return nil, fmt.Errorf("登录失败 (code: %s): %s", respBody.Code.Value.String(), respBody.Message.Value)
		}
		if !respBody.AccessKey.Valid || respBody.AccessKey.Value == "" {
			return nil, fmt.Errorf("登录失败 (code: %s): %s", respBody.Code.Value.String(), respBody.Message.Value)
		}
		uid := ""
		if respBody.Uid.Valid {
			uid = strconv.Itoa(respBody.Uid.Value)
		}
		return &SdkAccount{
			AccessKey: respBody.AccessKey.Value,
			Uid:       uid,
			Platform:  user.Platform,
			Channel:   user.Channel,
		}, nil
	}

	// 循环结束依然没成功，说明触发了兜底防线
	return nil, errors.New("触发过多验证码挑战，登录流程异常中断")
}
