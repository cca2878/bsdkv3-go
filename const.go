// Package bsdkv3 是 bilibili 游戏 SDK（external v3 登录协议）的 Go 客户端。
//
// 分层：物理网关 transport.Gateway ↔ interceptor 管道（retry/公共参数/时间戳/签名）
// ↔ service 业务层（登录/验证码）↔ validate 验证码求解（委托 gtrv-go）。
// 对外门面为 Client，业务入口为 Client.Auth。
package bsdkv3

// AppkeyPcr 是「公主连结 Re:Dive」国服的 bilibili SDK AppKey。
const AppkeyPcr = "fe8aac4e02f845b8ad67c427d48bfaf1"
