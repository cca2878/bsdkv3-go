// Package bsdkv3 provides input validation utilities for the BSdk V3 client.
//
// This file contains validation functions to ensure input data meets
// the requirements before making API calls.
package bsdkv3

import (
	"strings"
)

// ValidateUserInfo validates the UserInfo struct fields
func ValidateUserInfo(user UserInfo) error {
	if strings.TrimSpace(user.Username) == "" {
		return NewValidationError("Username", user.Username, "username cannot be empty")
	}
	
	if len(user.Username) > 100 {
		return NewValidationError("Username", user.Username, "username cannot exceed 100 characters")
	}
	
	if strings.TrimSpace(user.Password) == "" {
		return NewValidationError("Password", "[REDACTED]", "password cannot be empty")
	}
	
	if len(user.Password) > 200 {
		return NewValidationError("Password", "[REDACTED]", "password cannot exceed 200 characters")
	}
	
	// Platform and Channel are optional but if provided should not be empty strings
	if user.Platform != "" && strings.TrimSpace(user.Platform) == "" {
		return NewValidationError("Platform", user.Platform, "platform cannot be whitespace only")
	}
	
	if user.Channel != "" && strings.TrimSpace(user.Channel) == "" {
		return NewValidationError("Channel", user.Channel, "channel cannot be whitespace only")
	}
	
	return nil
}

// ValidateAppKey validates the application key
func ValidateAppKey(appKey string) error {
	if strings.TrimSpace(appKey) == "" {
		return NewValidationError("AppKey", "[REDACTED]", "app key cannot be empty")
	}
	
	if len(appKey) != 32 {
		return NewValidationError("AppKey", "[REDACTED]", "app key must be 32 characters long")
	}
	
	return nil
}