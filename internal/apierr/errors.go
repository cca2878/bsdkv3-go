// Package apierr 集中定义 SDK 的错误哨兵与结构化错误类型。
//
// 放在 internal 叶子包（仅依赖标准库）里，是为了让 service / validate / interceptor
// 各层都能引用而不产生 import 环；对外由根包 bsdkv3 重新导出，供调用方 errors.Is /
// errors.As 精确分流。
//
// 设计沿用哨兵 + 细分的惯例：细分错误用 %w 包住其所属的顶层哨兵，因此对细分做
// errors.Is 会同时命中它自己与顶层哨兵。
package apierr

import (
	"errors"
	"fmt"
)

// 顶层哨兵：调用方可用 errors.Is 归类失败大类。
var (
	// ErrInvalidRequest 请求在本地校验/编码阶段就不合法，尚未发出网络请求。
	ErrInvalidRequest = errors.New("bsdkv3: 请求参数非法")

	// ErrBootstrap 初始化引导阶段失败（拉取外部配置 / 登录 host 列表 / 登录密钥）。
	ErrBootstrap = errors.New("bsdkv3: 初始化引导失败")

	// ErrCipher 登录密钥相关失败（PEM/公钥解析、RSA 加密）。
	ErrCipher = errors.New("bsdkv3: 登录密钥处理失败")

	// ErrTransport 传输层失败（所有 host 均不可用等）。
	ErrTransport = errors.New("bsdkv3: 传输层失败")

	// ErrDecodeResponse 服务端响应反序列化失败。
	ErrDecodeResponse = errors.New("bsdkv3: 响应解析失败")

	// ErrCaptcha 人机验证（验证码）流程失败。
	ErrCaptcha = errors.New("bsdkv3: 人机验证失败")

	// ErrLogin 登录流程失败的顶层哨兵；LoginError 与下列登录细分都归属于它。
	ErrLogin = errors.New("bsdkv3: 登录失败")
)

// 细分哨兵：用 %w 包住顶层哨兵，errors.Is 会同时命中细分与顶层。
var (
	// ErrAllHostsFailed 高可用重试耗尽——所有 host 都失败。属于 ErrTransport。
	ErrAllHostsFailed = fmt.Errorf("%w: 所有 host 均失败", ErrTransport)

	// ErrTooManyCaptcha 触发过多验证码挑战，登录被兜底防线中断。属于 ErrLogin。
	ErrTooManyCaptcha = fmt.Errorf("%w: 触发过多验证码挑战", ErrLogin)

	// ErrMissingAccessKey 登录响应缺少 access_key。属于 ErrLogin。
	ErrMissingAccessKey = fmt.Errorf("%w: 响应缺少 access_key", ErrLogin)
)

// LoginError 承载服务端返回的业务失败码与消息，供调用方精确分流
// （例如区分密码错误与账号封禁）。errors.Is(err, ErrLogin) 恒为真，
// 也可 errors.As 取出 *LoginError 读 Code / Message。
type LoginError struct {
	Code    string // 服务端返回的业务 code（非 "0"）
	Message string // 服务端返回的提示消息
}

func (e *LoginError) Error() string {
	return fmt.Sprintf("bsdkv3: 登录失败 (code %s): %s", e.Code, e.Message)
}

// Unwrap 让 errors.Is(err, ErrLogin) 命中，把结构化错误归入登录大类。
func (e *LoginError) Unwrap() error { return ErrLogin }
