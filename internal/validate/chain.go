package validate

import (
	"context"
	"fmt"

	"github.com/cca2878/bsdkv3-go/internal/apierr"
)

// FailsafeChain 依然实现 Validator 接口
type FailsafeChain struct {
	validators []Validator
}

// NewFailsafeChain 构造降级链。传入的顺序极其重要：越靠前的优先级越高，越靠后的越兜底。
func NewFailsafeChain(v ...Validator) *FailsafeChain {
	return &FailsafeChain{
		validators: v,
	}
}

// Validate 核心降级逻辑：只要失败，就换下一个上！
func (c *FailsafeChain) Validate(ctx context.Context, challenge *ValidatorChallenge) (*ValidatorResult, error) {
	var lastErr error

	for _, v := range c.validators {
		// 尝试当前验证器
		result, err := v.Validate(ctx, challenge)

		if err == nil {
			// 只要有一个成功了，立刻中断链条并返回凭证
			return result, nil
		}

		// 记录失败原因（在实际工程中，这里应该打个 Warn 级别的日志，方便排查主节点为什么挂了）
		// logger.Warnf("第 %d 级验证器失败: %v，尝试降级...", i, err)
		lastErr = err

		// 检查 Context 是否已经被用户主动取消或超时，如果是，没必要再降级了，直接死掉
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	// 所有兜底方案全部阵亡
	return nil, fmt.Errorf("%w: 所有人机验证方案均失败, 最后错误: %w", apierr.ErrCaptcha, lastErr)
}
