// Package bsdkv3 provides custom error types for the BSdk V3 client.
//
// This file defines specific error types to help callers handle different
// types of failures appropriately.
package bsdkv3

import "fmt"

// Error types for different failure scenarios
var (
	// ErrInvalidCredentials indicates that the provided username/password combination is invalid
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")

	// ErrCaptchaFailed indicates that captcha validation failed
	ErrCaptchaFailed = fmt.Errorf("captcha validation failed")

	// ErrNetworkError indicates a network-related failure
	ErrNetworkError = fmt.Errorf("network error")

	// ErrConfigurationError indicates a client configuration error
	ErrConfigurationError = fmt.Errorf("configuration error")

	// ErrAPIError indicates an API-level error from the server
	ErrAPIError = fmt.Errorf("API error")

	// ErrAuthenticationFailed indicates general authentication failure
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
)

// ClientError represents a client-side error with additional context
type ClientError struct {
	Op      string // The operation that failed
	Err     error  // The underlying error
	Message string // Additional context message
}

func (e *ClientError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError creates a new ClientError with the given operation and underlying error
func NewClientError(op string, err error, message string) *ClientError {
	return &ClientError{
		Op:      op,
		Err:     err,
		Message: message,
	}
}

// APIError represents an error response from the API
type APIError struct {
	Code    string // Error code from the API
	Message string // Error message from the API
	Op      string // The operation that failed
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: API error (%s): %s", e.Op, e.Code, e.Message)
}

// NewAPIError creates a new APIError with the given details
func NewAPIError(op, code, message string) *APIError {
	return &APIError{
		Op:      op,
		Code:    code,
		Message: message,
	}
}

// ValidationError represents a validation-related error
type ValidationError struct {
	Field   string // The field that failed validation
	Value   string // The invalid value
	Message string // Description of what's wrong
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s' with value '%s': %s", e.Field, e.Value, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, value, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
