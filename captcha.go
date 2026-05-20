package bsdkv3

import (
	"context"
	"fmt"
	"net/http"
)

// 验证码API

const captchaAPIPath = "/api/client/start_captcha"

type captchaParams struct {
	CaptchaType string `form:"captcha_type"`
	SecCode     string `form:"seccode"`
	Validate    string `form:"validate"`
	GtUserId    string `form:"gt_user_id"`
	CToken      string `form:"ctoken"`
	Challenge   string `form:"challenge"`
}

type startCaptchaReq struct {
	baseReqParams
	Version string `form:"version"`
}

type startCaptchaResp struct {
	CaptchaType optionalValue[int]    `json:"captcha_type"`
	Gs          optionalValue[int]    `json:"gs"`
	Gt          optionalValue[string] `json:"gt"`
	Challenge   optionalValue[string] `json:"challenge"`
	GtUserId    optionalValue[string] `json:"gt_user_id"`
}

func (rq startCaptchaReq) validate() error {
	if rq.Version == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}

var startCaptchaAPI = endpoint[startCaptchaReq, startCaptchaResp]{
	Method: http.MethodPost,
	Path:   captchaAPIPath,
}

func (c *Client) handleCaptcha(ctx context.Context) (*captchaParams, error) {
	resp, err := execAPI(ctx, c, startCaptchaAPI, startCaptchaReq{
		baseReqParams: c.baseReqParams,
		Version:       defaultCaptchaVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("发送验证码请求失败: %w", err)
	}

	if resp.Raw.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("验证码请求返回非成功状态码: %d", resp.Raw.StatusCode())
	}

	ret, err := c.validator.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("验证码校验失败: %w", err)
	}
	if err := checkValidatorResult(ret); err != nil {
		return nil, fmt.Errorf("验证码校验结果无效: %w", err)
	}

	return &captchaParams{
		CaptchaType: defaultCaptchaType,
		Validate:    ret.Validate,
		Challenge:   ret.Challenge,
		GtUserId:    ret.GtUserId,
		SecCode:     ret.Validate + "|jordan",
		CToken:      ret.GtUserId,
	}, nil
}
