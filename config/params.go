package config

// BaseReqParams 各 API 请求的公共表单字段。
type BaseReqParams struct {
	CurBuvid          string `form:"cur_buvid"`
	OldBuvid          string `form:"old_buvid"`
	UdId              string `form:"udid"`
	BdId              string `form:"bd_id"`
	SdkType           string `form:"sdk_type"`
	VersionCode       string `form:"version_code"`
	MerchantId        string `form:"merchant_id"`
	ServerId          string `form:"server_id"`
	Version           string `form:"version"`
	DomainSwitchCount string `form:"domain_switch_count"`
	ApkSign           string `form:"apk_sign"`
	PlatformType      string `form:"platform_type"`
	AppVer            string `form:"app_ver"`
	SdkLogType        string `form:"sdk_log_type"`
	CurrentEnv        string `form:"current_env"`
	SdkVer            string `form:"sdk_ver"`
	AppId             string `form:"app_id"`
	Platform          string `form:"platform"`
	ChannelId         string `form:"channel_id"`
	GameId            string `form:"game_id"`
	Timestamp         string `form:"timestamp"`
	OriginalDomain    string `form:"original_domain"`
	Domain            string `form:"domain"`
	Sign              string `form:"sign"`
}

func NewDefaultBaseReqParams() BaseReqParams {
	return BaseReqParams{
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
}
