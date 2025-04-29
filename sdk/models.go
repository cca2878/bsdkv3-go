package sdk

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"bsdkv3/sdk/config"
)

//	type User interface {
//		GetUserInfo() UserInfo
//	}
type UserInfo struct {
	Username string
	Password string
}

func (u UserInfo) GetUserInfo() UserInfo {
	return u
}

// 请求相关常量
const (
	// MethodPost HTTP方法 - POST
	MethodPost = "POST"
)

// BaseRequest 基础请求公共字段
type BaseRequest struct {
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
}

// setDomainFromUrl 从 URL 更新 Domain 和 OriginalDomain 字段
func (b *BaseRequest) setDomainFromUrl(u *url.URL) {
	if u != nil {
		b.Domain = u.Host
		b.OriginalDomain = u.Scheme + "://" + u.Host
	}
}

// NewBaseRequest 创建基础请求对象（带默认值）
func NewBaseRequest() BaseRequest {
	// 从配置中获取默认值
	defaultReqConf := config.GetDefaultRequestConfig()

	return BaseRequest{
		CurBuvid:          defaultReqConf.CurBuvid,
		OldBuvid:          defaultReqConf.OldBuvid,
		UdId:              defaultReqConf.UdId,
		BdId:              defaultReqConf.BdId,
		SdkType:           defaultReqConf.SdkType,
		VersionCode:       defaultReqConf.VersionCode,
		MerchantId:        defaultReqConf.MerchantId,
		ServerId:          defaultReqConf.ServerId,
		Version:           defaultReqConf.Version,
		DomainSwitchCount: defaultReqConf.DomainSwitchCount,
		ApkSign:           defaultReqConf.ApkSign,
		PlatformType:      defaultReqConf.PlatformType,
		AppVer:            defaultReqConf.AppVer,
		SdkLogType:        defaultReqConf.SdkLogType,
		CurrentEnv:        defaultReqConf.CurrentEnv,
		SdkVer:            defaultReqConf.SdkVer,
		AppId:             defaultReqConf.AppId,
		Platform:          defaultReqConf.Platform,
		ChannelId:         defaultReqConf.ChannelId,
		GameId:            defaultReqConf.GameId,
		Timestamp:         strconv.FormatInt(time.Now().UnixMilli(), 10),
		Domain:            defaultReqConf.Domain,
		OriginalDomain:    defaultReqConf.OriginalDomain,
	}
}

// API路径常量

func parseModelUrl(host string, path string) (*url.URL, error) {
	u, err := url.Parse(host + path)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Req Resp base interface

type Request interface {

	// GetMethod Http Method
	GetMethod() string
	// GetUrl API Url
	GetUrl() (*url.URL, error)
}

type Response interface {
}

// BSdkV3ExtConf

const (
	// extConfAPIPath 获取配置API路径
	extConfAPIPath = "/api/external/config/v3"
)

type BSdkV3ExtConfReq struct {
	BaseRequest
}
type BSdkV3ExtConfResp struct {
	ConfigLoginHttps   string `json:"config_login_https"`
	ConfigAndroidHttps string `json:"config_login_android_https"`
}

func NewBSdkV3ExtConfReq() BSdkV3ExtConfReq {
	req := BSdkV3ExtConfReq{
		BaseRequest: NewBaseRequest(),
	}
	if reqUrl, err := req.GetUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return req
}

func (rq BSdkV3ExtConfReq) GetMethod() string {
	return MethodPost
}

func (rq BSdkV3ExtConfReq) GetUrl() (*url.URL, error) {
	return parseModelUrl(config.GetHostConfig().GetHost(config.HostTypeInitConf), extConfAPIPath)
}

// BSdkGetCipherV3

const (
	// getCipherV3APIPath 获取密钥API路径
	getCipherV3APIPath = "/api/external/issue/cipher/v3"
	cipherType         = "bili_login_rsa"
)

// BSdkGetCipherV3Req 获取密钥请求
type BSdkGetCipherV3Req struct {
	BaseRequest
	CipherType string `form:"cipher_type"`
}

// NewBSdkGetCipherV3Req 创建获取密钥请求实例
func NewBSdkGetCipherV3Req() BSdkGetCipherV3Req {
	req := BSdkGetCipherV3Req{
		BaseRequest: NewBaseRequest(),
		CipherType:  cipherType,
	}
	if reqUrl, err := req.GetUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return req
}

func (rq BSdkGetCipherV3Req) GetMethod() string {
	return MethodPost
}

func (rq BSdkGetCipherV3Req) GetUrl() (*url.URL, error) {
	return parseModelUrl(config.GetHostConfig().GetHost(config.HostTypeLoginHttps), getCipherV3APIPath)
}

// BSdkGetCipherV3Resp 获取密钥响应
type BSdkGetCipherV3Resp struct {
	CipherKey string `json:"cipher_key"`
	Hash      string `json:"hash"`
}

// BSdkV3Login

const (
	// loginAPIPath 登录API路径，与CaptLogin统一
	loginAPIPath = "/api/external/login/v3"
	BdInfo       = "cr_nmsl"
)

type BSdkV3LoginReq struct {
	BaseRequest
	BdInfo string `form:"bd_info"`
	UserId string `form:"user_id"`
	Pwd    string `form:"pwd"`
}

// NewBSdkV3LoginReq Factory Func
func NewBSdkV3LoginReq(user UserInfo) BSdkV3LoginReq {
	//userInfo := user.GetUserInfo()
	req := BSdkV3LoginReq{
		BaseRequest: NewBaseRequest(),
		BdInfo:      BdInfo,
		UserId:      user.Username,
		Pwd:         user.Password,
	}
	if reqUrl, err := req.GetUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return req
}

func (rq BSdkV3LoginReq) GetMethod() string {
	return config.MethodPost
}

func (rq BSdkV3LoginReq) GetUrl() (*url.URL, error) {
	return parseModelUrl(config.GetHostConfig().GetHost(config.HostTypeLoginHttps), loginAPIPath)
}

type BSdkV3LoginResp struct {
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

type CaptchaParams struct {
	CaptchaType string `form:"captcha_type"`
	SecCode     string `form:"seccode"`
	Validate    string `form:"validate"`
	GtUserId    string `form:"gt_user_id"`
	CToken      string `form:"ctoken"`
	Challenge   string `form:"challenge"`
}

func NewCaptchaParams(captParams CaptchaParams) CaptchaParams {
	if captParams.CaptchaType == "" {
		captParams.CaptchaType = config.DefaultCaptchaType
	}
	return captParams
}

type BSdkV3CaptLoginReq struct {
	BSdkV3LoginReq
	CaptchaParams
}

// NewBSdkV3CaptLoginReq Factory Func
func NewBSdkV3CaptLoginReq(user UserInfo, captParams CaptchaParams) BSdkV3CaptLoginReq {
	req := BSdkV3CaptLoginReq{
		BSdkV3LoginReq: NewBSdkV3LoginReq(user),
		CaptchaParams:  NewCaptchaParams(captParams),
	}
	if reqUrl, err := req.GetUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return req
}

func (rq BSdkV3CaptLoginReq) GetMethod() string {
	return config.MethodPost
}

func (rq BSdkV3CaptLoginReq) GetUrl() (*url.URL, error) {
	return parseModelUrl(config.GetHostConfig().GetHost(config.HostTypeLoginHttps), loginAPIPath)
}

type BSdkV3CaptLoginResp struct {
	BSdkV3LoginResp
}

// BSdkStartCaptcha

const (
	// captchaAPIPath 验证码API路径
	captchaAPIPath = "/api/client/start_captcha"
)

// BSdkStartCaptchaReq 开始验证码请求
type BSdkStartCaptchaReq struct {
	BaseRequest
	Version string `form:"version"`
}

// NewBSdkStartCaptchaReq 创建验证码请求实例
func NewBSdkStartCaptchaReq() BSdkStartCaptchaReq {
	req := BSdkStartCaptchaReq{
		BaseRequest: NewBaseRequest(),
		Version:     config.DefaultCaptchaVersion,
	}
	if reqUrl, err := req.GetUrl(); err == nil {
		req.setDomainFromUrl(reqUrl)
	}
	return req
}

func (rq BSdkStartCaptchaReq) GetMethod() string {
	return config.MethodPost
}

func (rq BSdkStartCaptchaReq) GetUrl() (*url.URL, error) {
	return parseModelUrl(config.GetHostConfig().GetHost(config.HostTypeLoginHttps), captchaAPIPath)
}

// BSdkStartCaptchaResp 开始验证码响应
type BSdkStartCaptchaResp struct {
	CaptchaType int    `json:"captcha_type"`
	Gs          int    `json:"gs"`
	Gt          string `json:"gt"`
	Challenge   string `json:"challenge"`
	GtUserId    string `json:"gt_user_id"`
}
