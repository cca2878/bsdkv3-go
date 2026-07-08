package validate

import (
	"context"
)

// 目前不需要传
type ValidatorChallenge struct {
	Type      string
	Challenge string
}

// ValidatorResult 定义验证结果的标准化结构
type ValidatorResult struct {
	Challenge string `json:"challenge"`
	Gt        string `json:"gt"`
	GtUserId  string `json:"gt_user_id"`
	Validate  string `json:"validate"`
}

// Validator 标准接口
type Validator interface {
	Validate(ctx context.Context, challenge *ValidatorChallenge) (*ValidatorResult, error)
}
