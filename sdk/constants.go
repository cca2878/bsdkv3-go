// Package bsdkv3 provides constants used throughout the BSdk V3 client library.
package bsdkv3

// API Version and SDK information
const (
	// SDKVersion represents the current version of this SDK
	SDKVersion = "1.0.0"

	// APIVersion represents the BSdk API version this client supports
	APIVersion = "3"

	// DefaultUserAgent is the default User-Agent header used in requests
	DefaultUserAgent = "Mozilla/5.0 BSGameSDK"

	// DefaultCVersion is the default client version header
	DefaultCVersion = "1"
)

// HTTP Headers
const (
	// ContentTypeForm is the content type for form data
	ContentTypeForm = "application/x-www-form-urlencoded"
)

// Response codes
const (
	// SuccessCode indicates a successful API response
	SuccessCode = "0"

	// CaptchaRequiredCode indicates captcha verification is required
	CaptchaRequiredCode = "1"
)
