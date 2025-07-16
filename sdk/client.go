// Package bsdkv3 provides a Go SDK client for BSdk V3 authentication and validation services.
//
// The SDK offers a simple API for user authentication with features including:
//   - Password encryption using RSA PKCS#1 v1.5
//   - Automatic captcha handling
//   - Configurable client options
//   - Comprehensive logging
//   - Host configuration management
//
// Basic usage:
//
//	client, err := bsdkv3.NewClient(config.AppkeyPcr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	user := bsdkv3.UserInfo{
//	    Username: "username",
//	    Password: "password",
//	}
//
//	account, err := client.Login(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The package supports various configuration options through the config package.
package bsdkv3

import (
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"bsdkv3-go/sdk/config"
	"bsdkv3-go/sdk/log"
)

// Client provides an HTTP client wrapper for BSdk V3 API operations.
// It handles authentication, request signing, and response processing.
//
// The client manages:
//   - HTTP client configuration and timeouts
//   - Request signing with MD5 signatures
//   - Password encryption using RSA public keys
//   - Automatic host configuration updates
//   - Context management for request cancellation
type Client struct {
	client      *resty.Client      // HTTP client for making requests
	formEncoder *form.Encoder      // Form encoder for request bodies
	publicKey   *rsa.PublicKey     // RSA public key for password encryption
	pwdHash     string             // Password hash salt from server
	appKey      string             // Application key for API authentication
	ctx         context.Context    // Client context for request lifecycle
	ctxCancel   context.CancelFunc // Function to cancel client context
	config      *config.Config     // Client configuration
}

// ClientOption defines a function type for configuring the Client
type ClientOption func(*Client)

// WithConfig sets a custom configuration for the client
func WithConfig(cfg *config.Config) ClientOption {
	return func(c *Client) {
		if cfg != nil {
			c.config = cfg
		}
	}
}

// NewClient creates a new Client instance with the given app key and options.
// It validates the app key, applies configuration options, and initializes
// the HTTP client with proper settings.
//
// The client automatically:
//   - Fetches external configuration from the API
//   - Updates host configurations
//   - Retrieves and parses the public key for encryption
//
// Parameters:
//   - appKey: The application key for API authentication (must be 32 characters)
//   - options: Optional configuration functions to customize the client
//
// Returns:
//   - *Client: A configured client ready for use
//   - error: Any error that occurred during initialization
func NewClient(appKey string, options ...ClientOption) (*Client, error) {
	// Validate the app key first
	if err := ValidateAppKey(appKey); err != nil {
		return nil, NewClientError("NewClient", err, "invalid app key provided")
	}

	// Create a context for the client lifetime
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize client with default configuration
	client := &Client{
		formEncoder: form.NewEncoder(),
		appKey:      appKey,
		ctx:         ctx,
		ctxCancel:   cancel,
		config:      config.NewDefaultConfig(),
	}

	// Apply any provided options
	for _, option := range options {
		option(client)
	}

	// Create and configure the HTTP client
	client.client = resty.New().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Mozilla/5.0 BSGameSDK",
			"cversion":     "1",
		}).
		SetTimeout(client.config.RequestTimeout)

	// Initialize configuration and cryptographic components
	if err := client.initialize(); err != nil {
		cancel() // Clean up the context
		return nil, NewClientError("NewClient", err, "failed to initialize client")
	}

	return client, nil
}

// initialize sets up the client by fetching configuration and cryptographic keys
func (c *Client) initialize() error {
	// Fetch external configuration
	confReq := newExtConfReq(c.config)
	var confResp extConfResp
	_, err := c.execReq(c.ctx, confReq, &confResp)
	if err != nil {
		return NewClientError("initialize", err, "failed to fetch external configuration")
	}

	if confResp.ConfigLoginHttps == "" {
		return NewClientError("initialize", ErrConfigurationError, "received empty login HTTPS configuration")
	}

	log.Info("更新登录Hosts")
	log.Debug("配置登录HTTPS: %s", confResp.ConfigLoginHttps)
	config.GetHostConfig().UpdateHosts(config.ParseHostsStr(config.HostTypeLoginHttps, confResp.ConfigLoginHttps))

	// Fetch cryptographic cipher
	cipherReq := newGetCipherV3Req(c.config)
	var cipherResp getCipherV3Resp
	_, err = c.execReq(c.ctx, cipherReq, &cipherResp)
	if err != nil {
		return NewClientError("initialize", err, "failed to fetch cipher key")
	}

	// Store password hash and parse public key
	c.pwdHash = cipherResp.Hash
	c.publicKey, err = parsePublicKeyFromPEM(cipherResp.CipherKey)
	if err != nil {
		return NewClientError("initialize", err, "failed to parse public key")
	}

	return nil
}

// calcSign computes the request signature using MD5 hash
func (c *Client) calcSign(requestBody interface{}) (string, error) {
	// Encode request body to form values
	values, err := c.formEncoder.Encode(requestBody)
	if err != nil {
		return "", NewClientError("calcSign", err, "failed to encode request body")
	}

	// Get and sort all keys
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Concatenate values in sorted key order
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(values.Get(k))
	}

	// Append app key
	sb.WriteString(c.appKey)

	// Calculate MD5 hash as signature
	sum := md5.Sum([]byte(sb.String()))
	return hex.EncodeToString(sum[:]), nil
}

// hashPwd encrypts the password using RSA public key
func (c *Client) hashPwd(pwd string) (string, error) {
	if c.publicKey == nil {
		return "", NewClientError("hashPwd", ErrConfigurationError, "public key not initialized")
	}
	
	if c.pwdHash == "" {
		return "", NewClientError("hashPwd", ErrConfigurationError, "password hash not initialized")
	}

	data, err := encryptPKCS1v15(c.publicKey, []byte(c.pwdHash+pwd))
	if err != nil {
		return "", NewClientError("hashPwd", err, "failed to encrypt password")
	}
	
	return base64.StdEncoding.EncodeToString(data), nil
}

// preReq prepares a request by adding signature and form data
func (c *Client) preReq(ctx context.Context, requestBody interface{}) (*resty.Request, error) {
	// Calculate signature
	sign, err := c.calcSign(requestBody)
	if err != nil {
		return nil, NewClientError("preReq", err, "failed to calculate signature")
	}

	// Encode request body as form values
	values, err := c.formEncoder.Encode(requestBody)
	if err != nil {
		return nil, NewClientError("preReq", err, "failed to encode request body")
	}

	// Add signature to form data
	values.Set("sign", sign)

	// Create and configure request
	req := c.client.R().
		SetContext(ctx).
		SetFormDataFromValues(values)

	return req, nil
}

// execReq executes a prepared request and handles the response
func (c *Client) execReq(ctx context.Context, request iRequest, result any) (*resty.Response, error) {
	req, err := c.preReq(ctx, request)
	if err != nil {
		log.Error("Failed to prepare request: %v", err)
		return nil, NewClientError("execReq", err, "request preparation failed")
	}

	url, err := request.getUrl()
	if err != nil {
		log.Error("Failed to get URL: %v", err)
		return nil, NewClientError("execReq", err, "failed to get request URL")
	}

	log.Debug("Sending request to: %s", url.String())

	// Send request and handle response
	resp, err := req.SetResult(result).Execute(request.getMethod(), url.String())
	if err != nil {
		log.Error("Request failed: %v", err)
		return resp, NewClientError("execReq", ErrNetworkError, 
			fmt.Sprintf("request to %s failed: %v", url.String(), err))
	}

	log.Debug("Received response: status=%d, length=%d", resp.StatusCode(), len(resp.Body()))

	if resp.StatusCode() != http.StatusOK {
		log.Error("Request failed with status code: %d", resp.StatusCode())
		return resp, NewClientError("execReq", ErrNetworkError, 
			fmt.Sprintf("request failed with status code: %d", resp.StatusCode()))
	}

	return resp, nil
}

// Login authenticates a user with the provided credentials and handles captcha if required.
//
// The login process:
//  1. Validates the user information
//  2. Encrypts the password using the public key
//  3. Attempts initial login
//  4. Handles captcha verification if required
//  5. Returns account information on success
//
// Parameters:
//   - user: UserInfo containing username, password, and optional platform/channel info
//
// Returns:
//   - *SdkAccount: Account information including UID and access key
//   - error: Any error that occurred during the login process
func (c *Client) Login(user UserInfo) (*SdkAccount, error) {
	// Validate user input
	if err := ValidateUserInfo(user); err != nil {
		return nil, NewClientError("Login", err, "invalid user information")
	}

	// Encrypt the password
	encryptedPassword, err := c.hashPwd(user.Password)
	if err != nil {
		return nil, NewClientError("Login", err, "failed to encrypt password")
	}

	// Create login request
	loginReq := newLoginReq(c.config).(*loginReq)
	loginReq.UserId = user.Username
	loginReq.Pwd = encryptedPassword

	// Attempt initial login
	var loginResp loginResp
	_, err = c.execReq(c.ctx, loginReq, &loginResp)
	if err != nil {
		return nil, NewClientError("Login", err, "login request failed")
	}

	// Check if captcha is required
	if loginResp.NeedCaptcha != nil && *loginResp.NeedCaptcha == "1" {
		return c.handleCaptchaLogin(user, encryptedPassword)
	}

	// Process normal login response
	return c.processLoginResponse(&loginResp, user, "Login")
}

// handleCaptchaLogin handles the captcha verification and login process
func (c *Client) handleCaptchaLogin(user UserInfo, encryptedPassword string) (*SdkAccount, error) {
	// Get captcha parameters
	captchaParams, err := c.handleCaptcha(c.ctx)
	if err != nil {
		return nil, NewClientError("handleCaptchaLogin", err, "captcha verification failed")
	}

	// Create captcha login request
	captLoginReq := newCaptLoginReq(c.config, *captchaParams).(*captLoginReq)
	captLoginReq.UserId = user.Username
	captLoginReq.Pwd = encryptedPassword

	// Execute captcha login request
	var captLoginResp captLoginResp
	_, err = c.execReq(c.ctx, captLoginReq, &captLoginResp)
	if err != nil {
		return nil, NewClientError("handleCaptchaLogin", err, "captcha login request failed")
	}

	// Process captcha login response
	return c.processLoginResponse(&captLoginResp.loginResp, user, "handleCaptchaLogin")
}

// processLoginResponse processes the login response and creates SdkAccount
func (c *Client) processLoginResponse(resp *loginResp, user UserInfo, operation string) (*SdkAccount, error) {
	if resp.Code == nil {
		return nil, NewClientError(operation, ErrAPIError, "missing response code")
	}

	// Check for successful login (code "0")
	if resp.Code.String() == "0" {
		if resp.AccessKey == nil {
			return nil, NewClientError(operation, ErrAPIError, "missing access key in successful response")
		}
		if resp.Uid == nil {
			return nil, NewClientError(operation, ErrAPIError, "missing UID in successful response")
		}

		return &SdkAccount{
			AccessKey: *resp.AccessKey,
			Uid:       strconv.Itoa(*resp.Uid),
			Platform:  user.Platform,
			Channel:   user.Channel,
		}, nil
	}

	// Handle error responses
	message := "unknown error"
	if resp.Message != nil {
		message = *resp.Message
	}

	// Create specific error based on response code
	apiErr := NewAPIError(operation, resp.Code.String(), message)
	
	// Map common error codes to specific error types
	switch resp.Code.String() {
	case "-629", "-626": // Common invalid credential codes
		return nil, NewClientError(operation, ErrInvalidCredentials, apiErr.Error())
	default:
		return nil, NewClientError(operation, ErrAuthenticationFailed, apiErr.Error())
	}
}
// handleCaptcha manages the captcha verification process
func (c *Client) handleCaptcha(ctx context.Context) (*captchaParams, error) {
	// Request captcha parameters
	captchaReq := newStartCaptchaReq(c.config)

	// Prepare and send captcha request
	req, err := c.preReq(ctx, captchaReq)
	if err != nil {
		return nil, NewClientError("handleCaptcha", err, "failed to prepare captcha request")
	}

	url, err := captchaReq.getUrl()
	if err != nil {
		return nil, NewClientError("handleCaptcha", err, "failed to get captcha URL")
	}

	var captchaResp startCaptchaResp
	resp, err := req.SetResult(&captchaResp).Post(url.String())
	if err != nil {
		return nil, NewClientError("handleCaptcha", err, "failed to send captcha request")
	}

	// Check HTTP status
	if resp.StatusCode() != http.StatusOK {
		return nil, NewClientError("handleCaptcha", ErrNetworkError, 
			fmt.Sprintf("captcha request returned status code: %d", resp.StatusCode()))
	}

	// Perform remote validation
	validator := NewRemoteValidator()
	validationResult, err := validator.Validate()
	if err != nil {
		return nil, NewClientError("handleCaptcha", ErrCaptchaFailed, 
			fmt.Sprintf("remote captcha validation failed: %v", err))
	}

	// Construct captcha parameters
	return &captchaParams{
		CaptchaType: "1",
		Validate:    validationResult.Validate,
		Challenge:   validationResult.Challenge,
		GtUserId:    validationResult.GtUserId,
		SecCode:     validationResult.Validate + "|jordan",
		CToken:      validationResult.GtUserId,
	}, nil
}

// startCaptcha retrieves captcha information (currently unused but kept for compatibility)
func (c *Client) startCaptcha(ctx context.Context) (*startCaptchaResp, error) {
	captchaReq := newStartCaptchaReq(c.config)

	var result startCaptchaResp
	_, err := c.execReq(ctx, captchaReq, &result)
	if err != nil {
		return nil, NewClientError("startCaptcha", err, "failed to get captcha information")
	}

	return &result, nil
}

// Close gracefully shuts down the client and cancels all pending requests.
// It should be called when the client is no longer needed to free resources.
func (c *Client) Close() {
	if c.ctxCancel != nil {
		log.Debug("Closing client and canceling pending requests")
		c.ctxCancel()
		c.ctxCancel = nil // Prevent double-close
	}
}
