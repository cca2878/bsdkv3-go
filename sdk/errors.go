package sdk

import "fmt"

// Error types for better error handling

// AuthError represents authentication-related errors
type AuthError struct {
	Code    string
	Message string
	Err     error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("authentication failed [%s]: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("authentication failed [%s]: %s", e.Code, e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new authentication error
func NewAuthError(code, message string, err error) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Component string
	Err       error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error in %s: %v", e.Component, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error
func NewConfigError(component string, err error) *ConfigError {
	return &ConfigError{
		Component: component,
		Err:       err,
	}
}

// CaptchaError represents captcha validation errors
type CaptchaError struct {
	Reason string
	Err    error
}

func (e *CaptchaError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("captcha validation failed: %s - %v", e.Reason, e.Err)
	}
	return fmt.Sprintf("captcha validation failed: %s", e.Reason)
}

func (e *CaptchaError) Unwrap() error {
	return e.Err
}

// NewCaptchaError creates a new captcha error
func NewCaptchaError(reason string, err error) *CaptchaError {
	return &CaptchaError{
		Reason: reason,
		Err:    err,
	}
}
