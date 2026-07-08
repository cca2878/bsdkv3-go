package service

import (
	"github.com/cca2878/bsdkv3-go/internal/interceptor"
	"github.com/cca2878/bsdkv3-go/internal/validate"
)

// Service 是 Layer 2 的业务载体对象
type Service struct {
	doer      interceptor.Invoker
	validator validate.Validator
	cipher    CipherResult
}

// NewService 强类型构造函数
func NewService(doer interceptor.Invoker, val validate.Validator, cipher CipherResult) *Service {
	return &Service{
		doer:      doer,
		validator: val,
		cipher:    cipher,
	}
}
