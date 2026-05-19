package bsdkv3

// baseReqParams 各 API 请求的公共表单字段。
type baseReqParams struct {
	// config: conf,

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

func newBaseReqParams(reqConf reqConf) baseReqParams {
	return baseReqParams{
		CurBuvid:          reqConf.CurBuvid,
		OldBuvid:          reqConf.OldBuvid,
		UdId:              reqConf.UdId,
		BdId:              reqConf.BdId,
		SdkType:           reqConf.SdkType,
		VersionCode:       reqConf.VersionCode,
		MerchantId:        reqConf.MerchantId,
		ServerId:          reqConf.ServerId,
		Version:           reqConf.Version,
		DomainSwitchCount: reqConf.DomainSwitchCount,
		ApkSign:           reqConf.ApkSign,
		PlatformType:      reqConf.PlatformType,
		AppVer:            reqConf.AppVer,
		SdkLogType:        reqConf.SdkLogType,
		CurrentEnv:        reqConf.CurrentEnv,
		SdkVer:            reqConf.SdkVer,
		AppId:             reqConf.AppId,
		Platform:          reqConf.Platform,
		ChannelId:         reqConf.ChannelId,
		GameId:            reqConf.GameId,
		Domain:            reqConf.Domain,
		OriginalDomain:    reqConf.OriginalDomain,
	}
}

// func modelURL(hosts *config.HostConfig, hostType config.HostType, path string) (*url.URL, error) {
// 	host := hosts.GetHost(hostType)
// 	u, err := url.Parse(host + path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return u, nil
// }
