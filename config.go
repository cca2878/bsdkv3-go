package bsdkv3

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

const (
	defaultCaptchaType    = "1"
	defaultCaptchaVersion = "1"
	defaultReqTimeout     = 5 // 秒
)

// reqConf 请求基础配置
type reqConf struct {
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

// clientConf 包含所有可配置项
type clientConf struct {
	TryTimes int
	Timeout  time.Duration
}

// NewClientConf 创建默认配置实例
func NewClientConf(options ...option[clientConf]) clientConf {
	clientconf := clientConf{
		Timeout:  defaultReqTimeout * time.Second,
		TryTimes: 3, // Default retry times
	}
	for _, opt := range options {
		opt.apply(&clientconf)
	}
	return clientconf
}

// NewReqConf 创建请求配置实例
func NewReqConf(options ...option[reqConf]) reqConf {
	reqconf := reqConf{
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
	}
	for _, opt := range options {
		opt.apply(&reqconf)
	}
	return reqconf
}

// Client配置选项

// WithConfigClient 设置 Client 配置
func WithConfigClient(conf clientConf) option[clientConf] {
	return optionFunc[clientConf](func(c *clientConf) {
		*c = conf
	})
}

// WithConfigClientReqTimeout 设置请求超时时间
func WithConfigClientReqTimeout(timeout time.Duration) option[clientConf] {
	return optionFunc[clientConf](func(c *clientConf) {
		c.Timeout = timeout
	})
}

// WithConfigClientReqTryTimes 设置请求重试次数
func WithConfigClientReqTryTimes(tryTimes int) option[clientConf] {
	return optionFunc[clientConf](func(c *clientConf) {
		c.TryTimes = tryTimes
	})
}

// 请求配置选项

// WithConfigReqConf 设置请求配置
func WithConfigReqConf(conf reqConf) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		*c = conf
	})
}

// WithConfigReqAppVer 设置 AppVer
func WithConfigReqAppVer(appVer string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.AppVer = appVer
	})
}

// WithConfigReqSdkVer 设置 SdkVer
func WithConfigReqSdkVer(sdkVer string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.SdkVer = sdkVer
	})
}

// WithConfigReqCurrentEnv 设置 CurrentEnv
func WithConfigReqCurrentEnv(currentEnv string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.CurrentEnv = currentEnv
	})
}

// WithConfigReqSdkLogType 设置 SdkLogType
func WithConfigReqSdkLogType(sdkLogType string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.SdkLogType = sdkLogType
	})
}

// WithConfigReqAppId 设置 AppId
func WithConfigReqAppId(appId string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.AppId = appId
	})
}

// WithConfigReqPlatform 设置 Platform
func WithConfigReqPlatform(platform string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.Platform = platform
	})
}

// WithConfigReqChannelId 设置 ChannelId
func WithConfigReqChannelId(channelId string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.ChannelId = channelId
	})
}

// WithConfigReqGameId 设置 GameId
func WithConfigReqGameId(gameId string) option[reqConf] {
	return optionFunc[reqConf](func(c *reqConf) {
		c.GameId = gameId
	})
}
