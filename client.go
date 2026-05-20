package bsdkv3

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/go-playground/form/v4"
	"resty.dev/v3"
)

type optionsBuilder struct {
	logger     Logger
	clientConf clientConf
	reqConf    reqConf
	validator  Validator
}

// Client 封装 HTTP 客户端和 API 调用。
type Client struct {
	// 基础设施
	clientConf clientConf
	hostMgr    *hostMgr
	logger     Logger

	// 业务基础
	httpClient  *resty.Client
	formEncoder *form.Encoder

	// 业务配置
	appKey        string
	publicKey     rsa.PublicKey
	baseReqParams baseReqParams

	pwdHash string

	validator Validator
}

func WithClientLogger(logger Logger) option[optionsBuilder] {
	return optionFunc[optionsBuilder](func(b *optionsBuilder) {
		b.logger = logger
	})
}
func WithClientConf(conf clientConf) option[optionsBuilder] {
	return optionFunc[optionsBuilder](func(b *optionsBuilder) {
		b.clientConf = conf
	})
}
func WithClientReqConf(conf reqConf) option[optionsBuilder] {
	return optionFunc[optionsBuilder](func(b *optionsBuilder) {
		b.reqConf = conf
	})
}

// WithValidator 设置验证码验证器。
func WithValidator(validator Validator) option[optionsBuilder] {
	return optionFunc[optionsBuilder](func(b *optionsBuilder) {
		b.validator = validator
	})
}

// NewClient 创建一个新的 Client 实例。
func NewClient(ctx context.Context, appKey string, options ...option[optionsBuilder]) (*Client, error) {

	builder := &optionsBuilder{
		logger:     discardLogger,
		clientConf: NewClientConf(),
		reqConf:    NewReqConf(),
	}
	for _, opt := range options {
		opt.apply(builder)
	}

	if builder.logger == nil {
		builder.logger = discardLogger
	}

	client := &Client{
		clientConf: builder.clientConf,
		hostMgr:    newHostManager([]string{defaultLoginHttpsHost}),
		logger:     builder.logger,
		httpClient: resty.New().
			SetHeaders(map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
				"User-Agent":   "Mozilla/5.0 BSGameSDK",
				"cversion":     "1",
			}).
			SetTimeout(builder.clientConf.Timeout),
		formEncoder:   form.NewEncoder(),
		appKey:        appKey,
		baseReqParams: newBaseReqParams(builder.reqConf),
	}
	remoteFallback := newRemoteValidator(
		resty.NewWithClient(client.httpClient.Client()),
		withValidatorLogger(client.logger),
	)
	client.validator = remoteFallback
	if builder.validator != nil {
		client.validator = newValidatorChain(builder.validator, remoteFallback)
	}

	if err := client.bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("配置失败: %w", err)
	}
	return client, nil
}
