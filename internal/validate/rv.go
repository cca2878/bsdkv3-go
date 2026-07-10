package validate

import (
	"context"

	"github.com/cca2878/gtrv-go"
)

// RemoteValidator 把 gtrv.Validator 适配为本包的 Validator 契约（补上 SDK 内部所需的
// challenge 形参；当前 gtrv 求解服务自管 geetest 会话、无需该形参）。
//
// 注意：验证码请求打向独立的 geetest 求解服务（非 bili 登录 API），因此其底层 HTTP
// 绝不能复用主 SDK 的业务管道（否则会被 commonParams/stamp/sign 污染并用错 appKey 签名）。
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

// NewRemoteValidator 用给定的 gtrv.Validator 构造适配器。
func NewRemoteValidator(v gtrv.Validator) Validator {
	return &RemoteValidator{validator: v}
}
