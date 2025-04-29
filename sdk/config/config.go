// Package config 提供项目中使用的所有常量配置
package config

import (
	"time"
)

// App Key
const (
	// AppkeyPcr 公主连结
	AppkeyPcr = "fe8aac4e02f845b8ad67c427d48bfaf1"
)

// 默认API Host
const (
	// defaultInitConfHost 默认初始Host
	defaultInitConfHost = "https://p.biligame.com"
	// defaultLoginHttpsHost 默认登录HTTPS Host
	defaultLoginHttpsHost = "https://line1-sdk-center-login-sh.biligame.net"
)

// 共有Form常量
const (
	CurBuvid          = "CR_NMSL"
	OldBuvid          = "CR_NMSL"
	UdId              = "CR_NMSL"
	BdId              = "cr-nmsl"
	SdkType           = "1"
	VersionCode       = "276"
	MerchantId        = "1"
	ServerId          = "1592"
	Version           = "3"
	DomainSwitchCount = "0"
	ApkSign           = "crnmsl"
	PlatformType      = "3"
	AppVer            = "8.1.0"
	SdkLogType        = "1"
	CurrentEnv        = "0"
	SdkVer            = "6.6.2"
	AppId             = "1370"
	Platform          = "3"
	ChannelId         = "1"
	GameId            = "1370"
	Domain            = "line1-sdk-center-login-sh.biligame.net"
	OriginalDomain    = "line1-sdk-center-login-sh.biligame.net"
)

// 验证码相关常量
const (
	// DefaultCaptchaType 默认验证码类型
	DefaultCaptchaType = "1"
	// DefaultCaptchaVersion 默认验证码版本
	DefaultCaptchaVersion = "1"
)

// 请求相关常量
const (
	// MethodPost HTTP方法 - POST
	MethodPost = "POST"

	// RequestTimeout 请求超时时间(秒)
	RequestTimeout = 5
)

// RequestConfig 请求基础配置
type RequestConfig struct {
	CurBuvid          string
	OldBuvid          string
	UdId              string
	BdId              string
	SdkType           string
	VersionCode       string
	MerchantId        string
	ServerId          string
	Version           string
	DomainSwitchCount string
	ApkSign           string
	PlatformType      string
	AppVer            string
	SdkLogType        string
	CurrentEnv        string
	SdkVer            string
	AppId             string
	Platform          string
	ChannelId         string
	GameId            string
	Domain            string
	OriginalDomain    string
}

// GetDefaultRequestConfig 获取默认请求配置
func GetDefaultRequestConfig() RequestConfig {
	return RequestConfig{
		CurBuvid:          CurBuvid,
		OldBuvid:          OldBuvid,
		UdId:              UdId,
		BdId:              BdId,
		SdkType:           SdkType,
		VersionCode:       VersionCode,
		MerchantId:        MerchantId,
		ServerId:          ServerId,
		Version:           Version,
		DomainSwitchCount: DomainSwitchCount,
		ApkSign:           ApkSign,
		PlatformType:      PlatformType,
		AppVer:            AppVer,
		SdkLogType:        SdkLogType,
		CurrentEnv:        CurrentEnv,
		SdkVer:            SdkVer,
		AppId:             AppId,
		Platform:          Platform,
		ChannelId:         ChannelId,
		GameId:            GameId,
		Domain:            Domain,
		OriginalDomain:    OriginalDomain,
	}
}

// GetRequestTimeout 获取请求超时时间
func GetRequestTimeout() time.Duration {
	return RequestTimeout * time.Second
}
