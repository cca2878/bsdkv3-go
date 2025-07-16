package sdk

import (
	"errors"
	"testing"
)

func TestAuthError(t *testing.T) {
	baseErr := errors.New("underlying error")

	// Test with underlying error
	authErr := NewAuthError("401", "Unauthorized", baseErr)
	expectedMsg := "authentication failed [401]: Unauthorized - underlying error"
	if authErr.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, authErr.Error())
	}

	if errors.Unwrap(authErr) != baseErr {
		t.Error("Unwrap should return the underlying error")
	}

	// Test without underlying error
	authErr2 := NewAuthError("400", "Bad Request", nil)
	expectedMsg2 := "authentication failed [400]: Bad Request"
	if authErr2.Error() != expectedMsg2 {
		t.Errorf("Expected error message %q, got %q", expectedMsg2, authErr2.Error())
	}

	if errors.Unwrap(authErr2) != nil {
		t.Error("Unwrap should return nil when no underlying error")
	}
}

func TestConfigError(t *testing.T) {
	baseErr := errors.New("config parsing failed")

	configErr := NewConfigError("host", baseErr)
	expectedMsg := "configuration error in host: config parsing failed"
	if configErr.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, configErr.Error())
	}

	if errors.Unwrap(configErr) != baseErr {
		t.Error("Unwrap should return the underlying error")
	}
}

func TestCaptchaError(t *testing.T) {
	baseErr := errors.New("network timeout")

	// Test with underlying error
	captchaErr := NewCaptchaError("timeout", baseErr)
	expectedMsg := "captcha validation failed: timeout - network timeout"
	if captchaErr.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, captchaErr.Error())
	}

	if errors.Unwrap(captchaErr) != baseErr {
		t.Error("Unwrap should return the underlying error")
	}

	// Test without underlying error
	captchaErr2 := NewCaptchaError("queue too long", nil)
	expectedMsg2 := "captcha validation failed: queue too long"
	if captchaErr2.Error() != expectedMsg2 {
		t.Errorf("Expected error message %q, got %q", expectedMsg2, captchaErr2.Error())
	}

	if errors.Unwrap(captchaErr2) != nil {
		t.Error("Unwrap should return nil when no underlying error")
	}
}
