package bsdkv3

import (
	"errors"
	"strings"
	"testing"
)

func TestClientError(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name        string
		op          string
		err         error
		message     string
		wantContain []string
	}{
		{
			name:        "error with message",
			op:          "TestOperation",
			err:         baseErr,
			message:     "additional context",
			wantContain: []string{"TestOperation", "additional context", "base error"},
		},
		{
			name:        "error without message",
			op:          "TestOperation",
			err:         baseErr,
			message:     "",
			wantContain: []string{"TestOperation", "base error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientErr := NewClientError(tt.op, tt.err, tt.message)
			errStr := clientErr.Error()

			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("ClientError.Error() = %v, want to contain %v", errStr, want)
				}
			}

			// Test Unwrap
			if unwrapped := clientErr.Unwrap(); unwrapped != tt.err {
				t.Errorf("ClientError.Unwrap() = %v, want %v", unwrapped, tt.err)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name        string
		op          string
		code        string
		message     string
		wantContain []string
	}{
		{
			name:        "API error with details",
			op:          "Login",
			code:        "400",
			message:     "Invalid credentials",
			wantContain: []string{"Login", "400", "Invalid credentials"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := NewAPIError(tt.op, tt.code, tt.message)
			errStr := apiErr.Error()

			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("APIError.Error() = %v, want to contain %v", errStr, want)
				}
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		value       string
		message     string
		wantContain []string
	}{
		{
			name:        "validation error with details",
			field:       "username",
			value:       "invalid_user",
			message:     "must not be empty",
			wantContain: []string{"username", "invalid_user", "must not be empty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationErr := NewValidationError(tt.field, tt.value, tt.message)
			errStr := validationErr.Error()

			for _, want := range tt.wantContain {
				if !strings.Contains(errStr, want) {
					t.Errorf("ValidationError.Error() = %v, want to contain %v", errStr, want)
				}
			}
		})
	}
}
