package bsdkv3

// UserInfo 登录用户凭据与渠道信息。
type UserInfo struct {
	Username string
	Password string

	Platform string
	Channel  string
}

// SdkAccount 登录成功后返回的账号信息。
type SdkAccount struct {
	Uid       string
	AccessKey string

	Platform string
	Channel  string
}
