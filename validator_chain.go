package bsdkv3

import (
	"context"
	"errors"
	"fmt"
)

type validatorChain struct {
	validators []Validator
}

// newValidatorChain 创建有序的验证码验证链。
func newValidatorChain(validators ...Validator) Validator {
	return &validatorChain{
		validators: append([]Validator(nil), validators...),
	}
}

func (c *validatorChain) Validate(ctx context.Context) (*ValidatorResult, error) {
	var errs []error
	for i, v := range c.validators {
		if v == nil {
			continue
		}
		result, err := v.Validate(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("validator[%d]: %w", i, err))
			continue
		}
		if err := checkValidatorResult(result); err != nil {
			errs = append(errs, fmt.Errorf("validator[%d]: %w", i, err))
			continue
		}
		return result, nil
	}

	if len(errs) == 0 {
		return nil, fmt.Errorf("validator chain is empty")
	}
	return nil, fmt.Errorf("all validators failed: %w", errors.Join(errs...))
}

func checkValidatorResult(r *ValidatorResult) error {
	if r == nil {
		return fmt.Errorf("nil result")
	}
	if r.Challenge == "" || r.Gt == "" || r.Validate == "" {
		return fmt.Errorf("validation response missing required fields")
	}
	return nil
}
