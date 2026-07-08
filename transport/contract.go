package transport

import (
	"context"
	"net/url"
)

type Request struct {
	Method string
	// 上层应该只传相对URL，中间层会负责拼接完整URL
	URL        *url.URL
	Headers    map[string]string
	FormParams url.Values

	Body []byte
}

type Response struct {
	StatusCode int
	Body       []byte
}

// Gateway 就是开放给高级用户的顶级 SPI
type Gateway interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}
