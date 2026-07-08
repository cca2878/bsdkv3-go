package transport

import (
	"context"
	"net/url"
)

type Request struct {
	Method string
	// URL 既可为相对 URL，也可为绝对 URL：
	//   - 相对 URL：交由 RetryInterceptor 在每次尝试时解析出高可用 host 并拼成完整 URL；
	//     业务 API（登录/验证码/密钥）走此路径。
	//   - 绝对 URL：直接按原样请求（ResolveReference 下绝对 URL 优先，不会被 host 覆盖）；
	//     用于不经过 RetryInterceptor 的场景，如 bootstrap 拉取初始 host 列表。
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
