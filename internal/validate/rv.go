package validate

import (
	"context"

	"github.com/cca2878/gtrv-go"
)

// HTTPClient 是远程验证器所需的最小 HTTP 能力（标准库 *http.Client 即满足）。
//
// 注意：验证码请求打向独立的 geetest 求解服务（非 bili 登录 API），因此绝不能复用
// 主 SDK 的业务管道（否则会被 commonParams/stamp/sign 污染并用错 appKey 签名）。
type HTTPClient = gtrv.HTTPClient

type RemoteValidator struct {
	validator gtrv.Validator
}

func (rv *RemoteValidator) Validate(ctx context.Context, challenge *ValidatorChallenge) (*ValidatorResult, error) {
	// 这里的 challenge 目前不需要传
	result, err := rv.validator.Validate(ctx)
	if err != nil {
		return nil, err
	}
	return &ValidatorResult{
		Challenge: result.Challenge,
		Gt:        result.Gt,
		GtUserId:  result.GtUserId,
		Validate:  result.Validate,
	}, nil
}

// NewRemoteValidator 用给定的 HTTP 客户端构造远程验证器。
func NewRemoteValidator(httpClient HTTPClient) Validator {
	return &RemoteValidator{validator: gtrv.NewRemoteValidator(httpClient)}
}
