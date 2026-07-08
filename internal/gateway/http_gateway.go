package gateway

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cca2878/bsdkv3-go/transport"
)

// httpGateway 是基于标准库 net/http 的默认物理底座实现。
type httpGateway struct {
	client *http.Client
}

// HTTPGatewayOptions 默认网关的可选配置。
type HTTPGatewayOptions struct {
	Timeout int    // 单次请求超时（秒），<=0 表示不设置（由 ctx 控制）
	Proxy   string // 代理地址，空表示不使用
	// Transport 注入共享的底层 *http.Transport（统一 proxy/TLS/连接池，配一次全局复用）。
	// 设置后 Proxy 选项失效——代理应由该 Transport 自带。
	Transport *http.Transport
}

func WithHTTPGatewayTimeout(timeout int) Option[HTTPGatewayOptions] {
	return optionFunc[HTTPGatewayOptions](func(o *HTTPGatewayOptions) {
		o.Timeout = timeout
	})
}

func WithHTTPGatewayProxy(proxy string) Option[HTTPGatewayOptions] {
	return optionFunc[HTTPGatewayOptions](func(o *HTTPGatewayOptions) {
		o.Proxy = proxy
	})
}

// WithHTTPGatewayTransport 注入共享底层 *http.Transport（与游戏 API / 资源下载等统一
// proxy/TLS/连接池）。注入后本网关只负责管超时，代理/TLS 交由该 Transport。
func WithHTTPGatewayTransport(rt *http.Transport) Option[HTTPGatewayOptions] {
	return optionFunc[HTTPGatewayOptions](func(o *HTTPGatewayOptions) {
		o.Transport = rt
	})
}

// 确保默认引擎完美实现了高级用户也可以实现的 SPI
var _ transport.Gateway = (*httpGateway)(nil)

// NewHTTPGateway 构造标准库 http 网关。
func NewHTTPGateway(opts ...Option[HTTPGatewayOptions]) transport.Gateway {
	options := &HTTPGatewayOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	client := &http.Client{}
	if options.Timeout > 0 {
		client.Timeout = time.Duration(options.Timeout) * time.Second
	}
	switch {
	case options.Transport != nil:
		// 注入共享 transport 优先；代理/TLS 由它自带，Proxy 选项忽略。
		client.Transport = options.Transport
	case options.Proxy != "":
		if proxyURL, err := url.Parse(options.Proxy); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}

	return &httpGateway{client: client}
}

func (g *httpGateway) Do(ctx context.Context, req *transport.Request) (*transport.Response, error) {
	var (
		body       io.Reader
		fallbackCT string
	)
	switch {
	case len(req.Body) > 0:
		// 有二进制载荷（JSON 字符串或文件流等）直接透传，Content-Type 由上层写在 Headers 里
		body = bytes.NewReader(req.Body)
	case len(req.FormParams) > 0:
		// 仅在无 Body 时降级为表单编码
		body = strings.NewReader(req.FormParams.Encode())
		fallbackCT = "application/x-www-form-urlencoded"
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), body)
	if err != nil {
		return nil, err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	// 上层未显式设置 Content-Type 时，为表单请求补上默认值
	if fallbackCT != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", fallbackCT)
	}

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &transport.Response{
		StatusCode: resp.StatusCode,
		Body:       data,
	}, nil
}
