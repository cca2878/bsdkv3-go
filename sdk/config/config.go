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

	// DefaultRequestTimeout 默认请求超时时间(秒)
	DefaultRequestTimeout = 5
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

// Config 包含所有可配置项
type Config struct {
	RequestConfig  RequestConfig
	RequestTimeout time.Duration
}

// NewDefaultConfig 创建默认配置实例
func NewDefaultConfig() *Config {
	return &Config{
		RequestConfig: RequestConfig{
			CurBuvid:          "CR_NMSL",
			OldBuvid:          "CR_NMSL",
			UdId:              "CR_NMSL",
			BdId:              "cr-nmsl",
			SdkType:           "1",
			VersionCode:       "276",
			MerchantId:        "1",
			ServerId:          "1592",
			Version:           "3",
			DomainSwitchCount: "0",
			ApkSign:           "crnmsl",
			PlatformType:      "3",
			AppVer:            "8.1.0",
			SdkLogType:        "1",
			CurrentEnv:        "0",
			SdkVer:            "6.6.2",
			AppId:             "1370",
			Platform:          "3",
			ChannelId:         "1",
			GameId:            "1370",
			Domain:            "line1-sdk-center-login-sh.biligame.net",
			OriginalDomain:    "line1-sdk-center-login-sh.biligame.net",
		},
		RequestTimeout: DefaultRequestTimeout * time.Second,
	}
}

// Option 是配置选项函数类型
type Option func(*Config)

// WithRequestConfig 设置整个请求配置
func WithRequestConfig(config RequestConfig) Option {
	return func(c *Config) {
		c.RequestConfig = config
	}
}

// WithRequestTimeout 设置请求超时时间
func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.RequestTimeout = timeout
	}
}

// WithCurBuvid 设置 CurBuvid
func WithCurBuvid(v string) Option {
	return func(c *Config) {
		c.RequestConfig.CurBuvid = v
	}
}

// WithOldBuvid 设置 OldBuvid
func WithOldBuvid(v string) Option {
	return func(c *Config) {
		c.RequestConfig.OldBuvid = v
	}
}

// WithUdId 设置 UdId
func WithUdId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.UdId = v
	}
}

// WithBdId 设置 BdId
func WithBdId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.BdId = v
	}
}

// WithSdkType 设置 SdkType
func WithSdkType(v string) Option {
	return func(c *Config) {
		c.RequestConfig.SdkType = v
	}
}

// WithVersionCode 设置 VersionCode
func WithVersionCode(v string) Option {
	return func(c *Config) {
		c.RequestConfig.VersionCode = v
	}
}

// WithMerchantId 设置 MerchantId
func WithMerchantId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.MerchantId = v
	}
}

// WithServerId 设置 ServerId
func WithServerId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.ServerId = v
	}
}

// WithVersion 设置 Version
func WithVersion(v string) Option {
	return func(c *Config) {
		c.RequestConfig.Version = v
	}
}

// WithDomainSwitchCount 设置 DomainSwitchCount
func WithDomainSwitchCount(v string) Option {
	return func(c *Config) {
		c.RequestConfig.DomainSwitchCount = v
	}
}

// WithApkSign 设置 ApkSign
func WithApkSign(v string) Option {
	return func(c *Config) {
		c.RequestConfig.ApkSign = v
	}
}

// WithPlatformType 设置 PlatformType
func WithPlatformType(v string) Option {
	return func(c *Config) {
		c.RequestConfig.PlatformType = v
	}
}

// WithAppVer 设置 AppVer
func WithAppVer(v string) Option {
	return func(c *Config) {
		c.RequestConfig.AppVer = v
	}
}

// WithSdkLogType 设置 SdkLogType
func WithSdkLogType(v string) Option {
	return func(c *Config) {
		c.RequestConfig.SdkLogType = v
	}
}

// WithCurrentEnv 设置 CurrentEnv
func WithCurrentEnv(v string) Option {
	return func(c *Config) {
		c.RequestConfig.CurrentEnv = v
	}
}

// WithSdkVer 设置 SdkVer
func WithSdkVer(v string) Option {
	return func(c *Config) {
		c.RequestConfig.SdkVer = v
	}
}

// WithAppId 设置 AppId
func WithAppId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.AppId = v
	}
}

// WithPlatform 设置 Platform
func WithPlatform(v string) Option {
	return func(c *Config) {
		c.RequestConfig.Platform = v
	}
}

// WithChannelId 设置 ChannelId
func WithChannelId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.ChannelId = v
	}
}

// WithGameId 设置 GameId
func WithGameId(v string) Option {
	return func(c *Config) {
		c.RequestConfig.GameId = v
	}
}

// WithDomain 设置 Domain
func WithDomain(v string) Option {
	return func(c *Config) {
		c.RequestConfig.Domain = v
	}
}

// WithOriginalDomain 设置 OriginalDomain
func WithOriginalDomain(v string) Option {
	return func(c *Config) {
		c.RequestConfig.OriginalDomain = v
	}
}

// NewConfig 使用选项创建自定义配置
func NewConfig(options ...Option) *Config {
	cfg := NewDefaultConfig()
	for _, option := range options {
		option(cfg)
	}
	return cfg
}
