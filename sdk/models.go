package bsdkv3

import (
	"bsdkv3-go/sdk/config"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//	type User interface {
//		GetUserInfo() UserInfo
//	}
type UserInfo struct {
	Username string
	Password string

	Platform string
	Channel  string
}

type SdkAccount struct {
	Uid       string
	AccessKey string

	Platform string
	Channel  string
}

// Req Resp base interface

type iRequest interface {
	setConfig(config *config.Config)
	// GetMethod Http Method
	getMethod() string
	// GetUrl API Url
	getUrl() (*url.URL, error)
}

type iResponse interface {
}

// baseRequest 基础请求公共字段
type baseRequest struct {
	// 关联的配置
	config *config.Config

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

// SetConfig 设置关联的配置
func (b *baseRequest) setConfig(config *config.Config) {
	b.config = config
}

// setDomainFromUrl 从 URL 更新 Domain 和 OriginalDomain 字段
func (b *baseRequest) setDomainFromUrl(u *url.URL) {
	if u != nil {
		b.Domain = u.Host
		b.OriginalDomain = u.Scheme + "://" + u.Host
	}
}

// newBaseRequest 创建基础请求对象
func newBaseRequest(conf *config.Config) baseRequest {
	// 确定使用哪个配置
	var reqConf config.RequestConfig
	if conf != nil {
		reqConf = conf.RequestConfig
	} else {
		reqConf = config.NewDefaultConfig().RequestConfig
	}

	return baseRequest{
		config: conf,

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
		Timestamp:         strconv.FormatInt(time.Now().UnixMilli(), 10),
		Domain:            reqConf.Domain,
		OriginalDomain:    reqConf.OriginalDomain,
	}
}

func parseModelUrl(hostType config.HostType, path string) (*url.URL, error) {
	host := config.GetHostConfig().GetHost(hostType)
	u, err := url.Parse(host + path)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// BSdkV3ExtConf

const (
	// extConfAPIPath 获取配置API路径
	extConfAPIPath = "/api/external/config/v3"
)

type extConfReq struct {
	baseRequest
}
type extConfResp struct {
	ConfigLoginHttps   string `json:"config_login_https"`
	ConfigAndroidHttps string `json:"config_login_android_https"`
}

func newExtConfReq(conf *config.Config) iRequest {
	req := extConfReq{
		baseRequest: newBaseRequest(conf),
	}
	if reqUrl, err := req.getUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return &req
}

func (rq extConfReq) getMethod() string {
	return http.MethodPost
}

func (rq extConfReq) getUrl() (*url.URL, error) {
	return parseModelUrl(config.HostTypeInitConf, extConfAPIPath)
}

// BSdkGetCipherV3

const (
	// getCipherV3APIPath 获取密钥API路径
	getCipherV3APIPath = "/api/external/issue/cipher/v3"
	cipherType         = "bili_login_rsa"
)

// getCipherV3Req 获取密钥请求
type getCipherV3Req struct {
	baseRequest
	CipherType string `form:"cipher_type"`
}

// newGetCipherV3Req 创建获取密钥请求实例
func newGetCipherV3Req(conf *config.Config) iRequest {
	req := getCipherV3Req{
		baseRequest: newBaseRequest(conf),
		CipherType:  cipherType,
	}
	if reqUrl, err := req.getUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return &req
}

func (rq getCipherV3Req) getMethod() string {
	return http.MethodPost
}

func (rq getCipherV3Req) getUrl() (*url.URL, error) {
	return parseModelUrl(config.HostTypeLoginHttps, getCipherV3APIPath)
}

// getCipherV3Resp 获取密钥响应
type getCipherV3Resp struct {
	CipherKey string `json:"cipher_key"`
	Hash      string `json:"hash"`
}

// BSdkV3Login

const (
	// loginAPIPath 登录API路径，与CaptLogin统一
	loginAPIPath = "/api/external/login/v3"
	bdInfo       = "cr_nmsl"
)

type loginReq struct {
	baseRequest
	BdInfo string `form:"bd_info"`
	UserId string `form:"user_id"`
	Pwd    string `form:"pwd"`
}

// newLoginReq Factory Func
func newLoginReq(conf *config.Config) iRequest {
	//userInfo := user.GetUserInfo()
	req := loginReq{
		baseRequest: newBaseRequest(conf),
		BdInfo:      bdInfo,
	}
	if reqUrl, err := req.getUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return &req
}

func (rq loginReq) getMethod() string {
	return http.MethodPost
}

func (rq loginReq) getUrl() (*url.URL, error) {
	return parseModelUrl(config.HostTypeLoginHttps, loginAPIPath)
}

type loginResp struct {
	// Code 你B程序员num和str混用
	Code *json.Number `json:"code"`

	NeedCaptcha      *string `json:"need_captch"`
	Nonce            *string `json:"nonce"`
	AccessKey        *string `json:"access_key"`
	Expires          *int    `json:"expires"`
	RealnameVerified *int    `json:"realname_verified"`
	Uid              *int    `json:"uid"`
	Uname            *string `json:"uname"`
	Message          *string `json:"message"`
}

// BSdkV3CaptLogin

type captchaParams struct {
	CaptchaType string `form:"captcha_type"`
	SecCode     string `form:"seccode"`
	Validate    string `form:"validate"`
	GtUserId    string `form:"gt_user_id"`
	CToken      string `form:"ctoken"`
	Challenge   string `form:"challenge"`
}

func newCaptchaParams(captParams captchaParams) captchaParams {
	if captParams.CaptchaType == "" {
		captParams.CaptchaType = config.DefaultCaptchaType
	}
	return captParams
}

type captLoginReq struct {
	loginReq
	captchaParams
}

// newCaptLoginReq Factory Func
func newCaptLoginReq(conf *config.Config, captParams captchaParams) iRequest {
	loginReq := newLoginReq(conf).(*loginReq)
	req := captLoginReq{
		loginReq:      *loginReq,
		captchaParams: newCaptchaParams(captParams),
	}
	if reqUrl, err := req.getUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return &req
}

func (rq captLoginReq) getMethod() string {
	return http.MethodPost
}

func (rq captLoginReq) getUrl() (*url.URL, error) {
	return parseModelUrl(config.HostTypeLoginHttps, loginAPIPath)
}

type captLoginResp struct {
	loginResp
}

// BSdkStartCaptcha

const (
	// captchaAPIPath 验证码API路径
	captchaAPIPath = "/api/client/start_captcha"
)

// startCaptchaReq 开始验证码请求
type startCaptchaReq struct {
	baseRequest
	Version string `form:"version"`
}

// newStartCaptchaReq 创建验证码请求实例
func newStartCaptchaReq(conf *config.Config) iRequest {
	req := startCaptchaReq{
		baseRequest: newBaseRequest(conf),
		Version:     config.DefaultCaptchaVersion,
	}
	if reqUrl, err := req.getUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return &req
}

func (rq startCaptchaReq) getMethod() string {
	return http.MethodPost
}

func (rq startCaptchaReq) getUrl() (*url.URL, error) {
	return parseModelUrl(config.HostTypeLoginHttps, captchaAPIPath)
}

// startCaptchaResp 开始验证码响应
type startCaptchaResp struct {
	CaptchaType int    `json:"captcha_type"`
	Gs          int    `json:"gs"`
	Gt          string `json:"gt"`
	Challenge   string `json:"challenge"`
	GtUserId    string `json:"gt_user_id"`
}
