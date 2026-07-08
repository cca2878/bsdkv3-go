package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cca2878/bsdkv3-go/internal/apierr"
	"github.com/cca2878/bsdkv3-go/internal/interceptor"
	"github.com/cca2878/bsdkv3-go/transport"
	"github.com/go-playground/form/v4"
)

type request interface {
	validate() error
}

type endpoint[Req request, Resp any] struct {
	Method string
	Path   string
}

func execAPI[Req request, Resp any](ctx context.Context, pipeDo interceptor.Invoker, endpoint endpoint[Req, Resp], reqBody Req) (*Resp, error) {
	if err := reqBody.validate(); err != nil {
		return nil, fmt.Errorf("请求参数校验失败: %w", err)
	}
	urlPath, err := url.Parse(endpoint.Path)
	if err != nil {
		return nil, fmt.Errorf("%w: 解析 URL 失败: %w", apierr.ErrInvalidRequest, err)
	}
	formValues, err := form.NewEncoder().Encode(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%w: 编码表单参数失败: %w", apierr.ErrInvalidRequest, err)
	}

	req := &transport.Request{
		Method:     endpoint.Method,
		URL:        urlPath,
		Headers:    make(map[string]string), // 预置非 nil map，供 StampInterceptor 写入请求头
		FormParams: formValues,
	}
	resp, err := pipeDo(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("%w: 管道返回空响应", apierr.ErrTransport)
	}
	// 非 2xx 一律视为基础设施/协议层失败，不再按业务响应解析。
	// （bili 登录 API 的业务失败走 HTTP 200 + code 字段；≥500 已在 RetryInterceptor
	// 内换 host 重试，此处主要拦截 4xx 以及重试耗尽后回落的响应。）
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%w: %d", apierr.ErrUnexpectedStatus, resp.StatusCode)
	}

	respBody := new(Resp)
	err = json.Unmarshal(resp.Body, respBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierr.ErrDecodeResponse, err)
	}
	return respBody, nil
}
