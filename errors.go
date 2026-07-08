package bsdkv3

import "github.com/cca2878/bsdkv3-go/internal/apierr"

// 错误哨兵与结构化类型在 internal/apierr 中定义（叶子包，避免 import 环），
// 这里于门面根包重新导出，供调用方 errors.Is / errors.As 精确分流失败原因。
var (
	// ErrInvalidRequest 请求在本地校验/编码阶段就不合法，尚未发出网络请求。
	ErrInvalidRequest = apierr.ErrInvalidRequest
	// ErrBootstrap 初始化引导阶段失败（外部配置 / 登录 host 列表 / 登录密钥）。
	ErrBootstrap = apierr.ErrBootstrap
	// ErrCipher 登录密钥相关失败（PEM/公钥解析、RSA 加密）。
	ErrCipher = apierr.ErrCipher
	// ErrTransport 传输层失败；ErrAllHostsFailed 归属于它。
	ErrTransport = apierr.ErrTransport
	// ErrAllHostsFailed 高可用重试耗尽，所有 host 均失败。
	ErrAllHostsFailed = apierr.ErrAllHostsFailed
	// ErrDecodeResponse 服务端响应反序列化失败。
	ErrDecodeResponse = apierr.ErrDecodeResponse
	// ErrCaptcha 人机验证（验证码）流程失败。
	ErrCaptcha = apierr.ErrCaptcha
	// ErrLogin 登录失败顶层哨兵；LoginError / ErrTooManyCaptcha / ErrMissingAccessKey 均归属于它。
	ErrLogin = apierr.ErrLogin
	// ErrTooManyCaptcha 触发过多验证码挑战，登录被兜底防线中断。
	ErrTooManyCaptcha = apierr.ErrTooManyCaptcha
	// ErrMissingAccessKey 登录响应缺少 access_key。
	ErrMissingAccessKey = apierr.ErrMissingAccessKey
)

// LoginError 承载服务端登录失败码与消息，可用 errors.As 取出读 Code / Message。
type LoginError = apierr.LoginError
