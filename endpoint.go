package bsdkv3

import "resty.dev/v3"

type request interface {
	validate() error
}

type endpoint[Req request, Resp any] struct {
	Method string
	Path   string
}

// apiResponse 包装了 API 响应的结果和原始响应对象。
type apiResponse[T any] struct {
	Body T
	Raw  *resty.Response
}
