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

	respBody := new(Resp)
	err = json.Unmarshal(resp.Body, respBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierr.ErrDecodeResponse, err)
	}
	return respBody, nil
}
