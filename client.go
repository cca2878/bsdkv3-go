package bsdkv3

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/cca2878/bsdkv3-go/config"
	"github.com/cca2878/bsdkv3-go/internal/base"
	"github.com/cca2878/bsdkv3-go/internal/gateway"
	"github.com/cca2878/bsdkv3-go/internal/interceptor"
	"github.com/cca2878/bsdkv3-go/internal/service"
	"github.com/cca2878/bsdkv3-go/internal/validate"
	"github.com/cca2878/bsdkv3-go/transport"
	"github.com/cca2878/gtrv-go"
)

// clientConfig 外部用户传入的显式初始化配置
type clientConfig struct {
	BaseParams config.BaseReqParams // 全局公共参数
	Gateway    transport.Gateway    // 可选的物理底座网关
	Timeout    time.Duration        // 可选的全局超时设置(仅不使用默认网关时生效)
	RetryTimes int
	// Transport 共享的底层 *http.Transport（统一 proxy/TLS/连接池）。注入后，默认网关
	// 与默认验证码 HTTP 客户端都会复用它——集成方只需配一次 transport 即可全链路共享。
	// 若显式传入 WithClientGateway / WithClientCaptchaHTTPClient，则各自以显式值为准。
	Transport *http.Transport
	// CaptchaHTTPClient 验证码求解走的独立 HTTP 客户端（打向 geetest 求解服务，
	// 不经过 bili 登录业务管道）。为 nil 时使用带全局超时的标准库客户端。
	CaptchaHTTPClient gtrv.HTTPClient
}

func WithClientGateway(gw transport.Gateway) Option[clientConfig] {
	return optionFunc[clientConfig](func(c *clientConfig) {
		c.Gateway = gw
	})
}

// WithClientTransport 注入共享的底层 *http.Transport，默认网关与默认验证码客户端都会复用它
// （与 go-autopcr 游戏 API / 资源下载统一 proxy/TLS/连接池）。
func WithClientTransport(rt *http.Transport) Option[clientConfig] {
	return optionFunc[clientConfig](func(c *clientConfig) {
		c.Transport = rt
	})
}

// WithClientCaptchaHTTPClient 注入验证码求解专用的 HTTP 客户端（如需代理或用于测试打桩）。
func WithClientCaptchaHTTPClient(hc gtrv.HTTPClient) Option[clientConfig] {
	return optionFunc[clientConfig](func(c *clientConfig) {
		c.CaptchaHTTPClient = hc
	})
}

func WithClientBaseParams(params config.BaseReqParams) Option[clientConfig] {
	return optionFunc[clientConfig](func(c *clientConfig) {
		c.BaseParams = params
	})
}

func WithClientRetryTimes(times int) Option[clientConfig] {
	return optionFunc[clientConfig](func(c *clientConfig) {
		c.RetryTimes = times
	})
}

// Client SDK 的唯一对外门面
type Client struct {
	// 【内部基建】：动态可变的管道原子指针，对外隐藏
	pipeline atomic.Pointer[interceptor.Pipeline]
	// 【内部配置】
	appKey string
	config clientConfig
	// 【业务门面】：外部开发者唯一能点出来的业务入口
	Auth *service.Service
}

// NewClient 整个 SDK 的总装车间（采用 Fail-Fast 隐式预加载模式）
func NewClient(ctx context.Context, appKey string, opts ...Option[clientConfig]) (*Client, error) {
	// 预设配置
	conf := clientConfig{
		Gateway:    nil,
		BaseParams: config.NewDefaultBaseReqParams(),
		Timeout:    20 * time.Second,
		RetryTimes: 3,
	}
	for _, opt := range opts {
		opt.apply(&conf)
	}

	// 制造物理底座网关
	// 可选配置项通过 Option 模式注入
	if conf.Gateway == nil {
		gwOpts := []gateway.Option[gateway.HTTPGatewayOptions]{
			gateway.WithHTTPGatewayTimeout(int(conf.Timeout.Seconds())),
		}
		if conf.Transport != nil {
			gwOpts = append(gwOpts, gateway.WithHTTPGatewayTransport(conf.Transport))
		}
		conf.Gateway = gateway.NewHTTPGateway(gwOpts...)
	}
	// 组装“第一级火箭”（用于拉取初始配置的临时轻量管道）
	initPipe := interceptor.NewPipeline(conf.Gateway).
		Use(interceptor.NewCommonParamsInterceptor(conf.BaseParams)).
		Use(interceptor.NewStampInterceptor()).
		Use(interceptor.NewSignInterceptor(appKey)) // 签名

	// 发射第一级火箭（调用解耦的纯函数拉取核心配置）
	hosts, err := service.FetchBootstrapHosts(ctx, initPipe.Do)
	if err != nil {
		return nil, err // 尽早失败，不暴露半成品对象
	}

	// 拿到配置，开始组装“第二级火箭”（终极业务管道）
	fullPipe := interceptor.NewPipeline(conf.Gateway).
		Use(interceptor.NewRetryInterceptor(conf.RetryTimes, base.NewHostManager(hosts))).
		Use(interceptor.NewCommonParamsInterceptor(conf.BaseParams)).
		Use(interceptor.NewStampInterceptor()).
		Use(interceptor.NewSignInterceptor(appKey))

	// 实例化主体
	client := &Client{
		config: conf,   // 保存配置
		appKey: appKey, // 保存 AppKey
	}
	client.pipeline.Store(fullPipe) // 将终极管道装入原子指针

	// 拉取登录密钥（RSA 公钥 + hash 盐）。走 fullPipe，直接享受高可用 host 路由与重试。
	cipher, err := service.FetchCipher(ctx, client.do)
	if err != nil {
		return nil, err // 尽早失败，不暴露半成品对象
	}

	// 组装验证码 Failsafe 降级链。
	// 验证码打向独立的 geetest 求解服务，必须用独立 HTTP 客户端，绝不能复用业务管道
	// （否则会被 commonParams/stamp/sign 污染并用错 appKey 签名）。
	captchaHTTP := conf.CaptchaHTTPClient
	if captchaHTTP == nil {
		hc := &http.Client{Timeout: conf.Timeout}
		if conf.Transport != nil {
			hc.Transport = conf.Transport // 与业务网关共享同一底层 transport
		}
		captchaHTTP = hc
	}
	remoteSolver := validate.NewRemoteValidator(captchaHTTP)
	failsafeValidator := validate.NewFailsafeChain(remoteSolver)

	// 装配 Layer 2 业务大脑，完成“闭包指针绑定”与“凭证分流”
	client.Auth = service.NewService(
		client.do, // 隐式捕获 client 指针，奠定热更新基础
		failsafeValidator,
		*cipher, // 密码加密凭证（RSA 公钥 + hash 盐）注入
	)

	return client, nil
}

// Do 门面代理方法（实现 transport.Invoker 签名）
// 彻底将 atomic.Load 的复杂性拦截在门面内部，让 Layer 2 彻底解耦
func (c *Client) do(ctx context.Context, req *transport.Request) (*transport.Response, error) {
	return c.pipeline.Load().Do(ctx, req)
}
