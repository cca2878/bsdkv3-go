package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cca2878/bsdkv3-go/internal/apierr"
	"github.com/cca2878/bsdkv3-go/internal/base"
	"github.com/cca2878/bsdkv3-go/internal/validate"
)

// 验证码API

const captchaAPIPath = "/api/client/start_captcha"

const (
	defaultCaptchaType    = "1"
	defaultCaptchaVersion = "1"
)

type captchaResult struct {
	CaptchaType string
	SecCode     string
	Validate    string
	GtUserId    string
	CToken      string
	Challenge   string
}

type startCaptchaReq struct {
	Version string `form:"version"`
}

type startCaptchaResp struct {
	CaptchaType base.OptionalValue[int]    `json:"captcha_type"`
	Gs          base.OptionalValue[int]    `json:"gs"`
	Gt          base.OptionalValue[string] `json:"gt"`
	Challenge   base.OptionalValue[string] `json:"challenge"`
	GtUserId    base.OptionalValue[string] `json:"gt_user_id"`
}

func (rq startCaptchaReq) validate() error {
	if rq.Version == "" {
		return fmt.Errorf("%w: version 不能为空", apierr.ErrInvalidRequest)
	}
	return nil
}

var startCaptchaAPI = endpoint[startCaptchaReq, startCaptchaResp]{
	Method: http.MethodPost,
	Path:   captchaAPIPath,
}

func (s *Service) handleCaptcha(ctx context.Context) (*captchaResult, error) {
	respBody, err := execAPI(ctx, s.doer, startCaptchaAPI, startCaptchaReq{
		Version: defaultCaptchaVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: 发送验证码请求失败: %w", apierr.ErrCaptcha, err)
	}

	ret, err := s.validator.Validate(ctx, &validate.ValidatorChallenge{
		Type:      defaultCaptchaType,
		Challenge: respBody.Challenge.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: 验证码校验失败: %w", apierr.ErrCaptcha, err)
	}

	return &captchaResult{
		CaptchaType: defaultCaptchaType,
		Validate:    ret.Validate,
		Challenge:   ret.Challenge,
		GtUserId:    ret.GtUserId,
		SecCode:     ret.Validate + "|jordan",
		CToken:      ret.GtUserId,
	}, nil
}
