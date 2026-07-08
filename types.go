package bsdkv3

import "github.com/cca2878/bsdkv3-go/internal/service"

// 业务类型在 internal/service 中定义，这里于门面根包重新导出，
// 使外部调用方能够命名 Login 的入参/出参类型（否则 internal 包不可被外部 import）。

// UserInfo 登录用户凭据与渠道信息。
type UserInfo = service.UserInfo

// SdkAccount 登录成功后返回的账号信息。
type SdkAccount = service.SdkAccount
