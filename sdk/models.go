// Package bsdkv3 provides model definitions and request/response structures for the BSdk V3 API.
//
// This file contains all the data structures used for API communication including:
//   - Request and response models
//   - User information structures
//   - Account information structures
//   - Base request/response interfaces
package bsdkv3

import (
	"bsdkv3-go/sdk/config"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// UserInfo contains user credentials and platform information for authentication.
//
// All fields should be provided for complete authentication, though Platform
// and Channel may be optional depending on the specific use case.
type UserInfo struct {
	Username string // User's login username (required)
	Password string // User's login password (required)
	Platform string // Platform identifier (optional)
	Channel  string // Channel identifier (optional)
}

// SdkAccount represents the authenticated user account information returned after successful login.
//
// This structure contains the essential information needed to make authenticated
// API calls on behalf of the user.
type SdkAccount struct {
	Uid       string // Unique user identifier
	AccessKey string // Access token for authenticated requests
	Platform  string // Platform the user authenticated from
	Channel   string // Channel the user authenticated through
}

// Request and Response Interfaces

// iRequest defines the interface that all API request types must implement.
// It provides methods for configuration, HTTP method determination, and URL construction.
type iRequest interface {
	setConfig(config *config.Config) // Set the configuration for this request
	getMethod() string               // Get the HTTP method (GET, POST, etc.)
	getUrl() (*url.URL, error)       // Get the complete URL for this request
}

// iResponse defines the interface for API response types.
// Currently a marker interface for type safety and future extensibility.
type iResponse interface {
}

// baseRequest contains common fields shared by all API requests.
// It includes device information, SDK version, platform details, and other
// metadata required by the BSdk V3 API.
type baseRequest struct {
	config *config.Config // Associated configuration

	// Device and SDK Information
	CurBuvid    string `form:"cur_buvid"`    // Current browser/device unique ID
	OldBuvid    string `form:"old_buvid"`    // Previous browser/device unique ID
	UdId        string `form:"udid"`         // Unique device identifier
	BdId        string `form:"bd_id"`        // Board/device identifier
	SdkType     string `form:"sdk_type"`     // SDK type identifier
	VersionCode string `form:"version_code"` // Application version code
	SdkVer      string `form:"sdk_ver"`      // SDK version
	AppVer      string `form:"app_ver"`      // Application version

	// Application Information
	MerchantId string `form:"merchant_id"` // Merchant identifier
	ServerId   string `form:"server_id"`   // Server identifier
	AppId      string `form:"app_id"`      // Application identifier
	GameId     string `form:"game_id"`     // Game identifier

	// Platform and Environment
	Platform     string `form:"platform"`      // Platform identifier
	PlatformType string `form:"platform_type"` // Platform type
	ChannelId    string `form:"channel_id"`    // Channel identifier
	CurrentEnv   string `form:"current_env"`   // Current environment

	// Request Metadata
	Version           string `form:"version"`             // API version
	Timestamp         string `form:"timestamp"`           // Request timestamp
	DomainSwitchCount string `form:"domain_switch_count"` // Domain switch count
	Domain            string `form:"domain"`              // Request domain
	OriginalDomain    string `form:"original_domain"`     // Original domain

	// Security and Signing
	ApkSign    string `form:"apk_sign"`     // APK signature
	SdkLogType string `form:"sdk_log_type"` // SDK log type
	Sign       string `form:"sign"`         // Request signature (added by client)
}

// setConfig associates a configuration with this request
func (b *baseRequest) setConfig(config *config.Config) {
	b.config = config
}

// setDomainFromUrl updates the Domain and OriginalDomain fields based on the provided URL
func (b *baseRequest) setDomainFromUrl(u *url.URL) {
	if u != nil {
		b.Domain = u.Host
		b.OriginalDomain = u.Scheme + "://" + u.Host
	}
}

// newBaseRequest creates a new baseRequest instance with values from the provided configuration.
// If no configuration is provided, it uses default values.
func newBaseRequest(conf *config.Config) baseRequest {
	// Determine which configuration to use
	var reqConf config.RequestConfig
	if conf != nil {
		reqConf = conf.RequestConfig
	} else {
		reqConf = config.NewDefaultConfig().RequestConfig
	}

	return baseRequest{
		config: conf,

		// Device and SDK Information
		CurBuvid:    reqConf.CurBuvid,
		OldBuvid:    reqConf.OldBuvid,
		UdId:        reqConf.UdId,
		BdId:        reqConf.BdId,
		SdkType:     reqConf.SdkType,
		VersionCode: reqConf.VersionCode,
		SdkVer:      reqConf.SdkVer,
		AppVer:      reqConf.AppVer,

		// Application Information
		MerchantId: reqConf.MerchantId,
		ServerId:   reqConf.ServerId,
		AppId:      reqConf.AppId,
		GameId:     reqConf.GameId,

		// Platform and Environment
		Platform:     reqConf.Platform,
		PlatformType: reqConf.PlatformType,
		ChannelId:    reqConf.ChannelId,
		CurrentEnv:   reqConf.CurrentEnv,

		// Request Metadata
		Version:           reqConf.Version,
		Timestamp:         strconv.FormatInt(time.Now().UnixMilli(), 10),
		DomainSwitchCount: reqConf.DomainSwitchCount,
		Domain:            reqConf.Domain,
		OriginalDomain:    reqConf.OriginalDomain,

		// Security and Signing
		ApkSign:    reqConf.ApkSign,
		SdkLogType: reqConf.SdkLogType,
	}
}

// parseModelUrl constructs a complete URL by combining a host from configuration with the given path
func parseModelUrl(hostType config.HostType, path string) (*url.URL, error) {
	host := config.GetHostConfig().GetHost(hostType)
	u, err := url.Parse(host + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL with host %s and path %s: %w", host, path, err)
	}
	return u, nil
}

// External Configuration API
//
// The external configuration API retrieves host configuration and other
// settings required for proper SDK operation.

const (
	extConfAPIPath = "/api/external/config/v3" // API path for external configuration
)

// extConfReq represents a request to fetch external configuration
type extConfReq struct {
	baseRequest
}

// extConfResp represents the response from the external configuration API
type extConfResp struct {
	ConfigLoginHttps   string `json:"config_login_https"`         // HTTPS login host configuration
	ConfigAndroidHttps string `json:"config_login_android_https"` // Android HTTPS login host configuration
}

// newExtConfReq creates a new external configuration request
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

// Cipher Key API
//
// The cipher key API retrieves the RSA public key and password hash
// used for encrypting user passwords during authentication.

const (
	getCipherV3APIPath = "/api/external/issue/cipher/v3" // API path for cipher key retrieval
	cipherType         = "bili_login_rsa"                // RSA cipher type for login
)

// getCipherV3Req represents a request to fetch the cipher key
type getCipherV3Req struct {
	baseRequest
	CipherType string `form:"cipher_type"` // Type of cipher requested
}

// newGetCipherV3Req creates a new cipher key request
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

// getCipherV3Resp represents the response from the cipher key API
type getCipherV3Resp struct {
	CipherKey string `json:"cipher_key"` // RSA public key in PEM format
	Hash      string `json:"hash"`       // Password hash salt
}

// User Login API
//
// The login API handles user authentication with username and password.
// It supports both normal login and captcha-protected login flows.

const (
	loginAPIPath = "/api/external/login/v3" // API path for login (shared with captcha login)
	bdInfo       = "cr_nmsl"                // Board info identifier
)

// loginReq represents a standard login request
type loginReq struct {
	baseRequest
	BdInfo string `form:"bd_info"` // Board info identifier
	UserId string `form:"user_id"` // Username for authentication
	Pwd    string `form:"pwd"`     // Encrypted password
}

// newLoginReq creates a new login request
func newLoginReq(conf *config.Config) iRequest {
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

// loginResp represents the response from a login request
type loginResp struct {
	Code *json.Number `json:"code"` // Response code (0 = success)

	// Login result fields
	NeedCaptcha      *string `json:"need_captch"`       // Whether captcha is required ("1" = yes)
	Nonce            *string `json:"nonce"`             // Security nonce
	AccessKey        *string `json:"access_key"`        // Access token for authenticated requests
	Expires          *int    `json:"expires"`           // Token expiration time
	RealnameVerified *int    `json:"realname_verified"` // Real name verification status
	Uid              *int    `json:"uid"`               // User ID
	Uname            *string `json:"uname"`             // Username
	Message          *string `json:"message"`           // Error or status message
}

// Captcha-Protected Login API
//
// When login requires captcha verification, additional parameters must be
// included in the login request.

// captchaParams contains the parameters required for captcha verification
type captchaParams struct {
	CaptchaType string `form:"captcha_type"` // Type of captcha (usually "1")
	SecCode     string `form:"seccode"`      // Security code from captcha validation
	Validate    string `form:"validate"`     // Validation token from captcha service
	GtUserId    string `form:"gt_user_id"`   // User ID from captcha service
	CToken      string `form:"ctoken"`       // Captcha token
	Challenge   string `form:"challenge"`    // Challenge string from captcha service
}

// newCaptchaParams creates captcha parameters with default values
func newCaptchaParams(captParams captchaParams) captchaParams {
	if captParams.CaptchaType == "" {
		captParams.CaptchaType = config.DefaultCaptchaType
	}
	return captParams
}

// captLoginReq represents a login request with captcha verification
type captLoginReq struct {
	loginReq      // Embed standard login request
	captchaParams // Add captcha-specific parameters
}

// newCaptLoginReq creates a new captcha login request
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

// captLoginResp represents the response from a captcha-protected login
type captLoginResp struct {
	loginResp // Embed standard login response
}

// Captcha Initialization API
//
// The captcha initialization API starts the captcha verification process
// and returns the necessary parameters for the captcha service.

const (
	captchaAPIPath = "/api/client/start_captcha" // API path for captcha initialization
)

// startCaptchaReq represents a request to initialize captcha verification
type startCaptchaReq struct {
	baseRequest
	Version string `form:"version"` // Captcha version
}

// newStartCaptchaReq creates a new captcha initialization request
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

// startCaptchaResp represents the response from captcha initialization
type startCaptchaResp struct {
	CaptchaType int    `json:"captcha_type"` // Type of captcha required
	Gs          int    `json:"gs"`           // Captcha service identifier
	Gt          string `json:"gt"`           // Captcha GT parameter
	Challenge   string `json:"challenge"`    // Challenge string for captcha
	GtUserId    string `json:"gt_user_id"`   // User ID for captcha service
}
