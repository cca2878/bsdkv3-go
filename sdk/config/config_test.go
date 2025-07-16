package config

import (
	"testing"
	"time"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	// Test default values
	if cfg.RequestTimeout != DefaultRequestTimeout*time.Second {
		t.Errorf("Expected timeout %v, got %v", DefaultRequestTimeout*time.Second, cfg.RequestTimeout)
	}

	if cfg.RequestConfig.SdkType != "1" {
		t.Errorf("Expected SdkType '1', got %s", cfg.RequestConfig.SdkType)
	}

	if cfg.RequestConfig.Platform != "3" {
		t.Errorf("Expected Platform '3', got %s", cfg.RequestConfig.Platform)
	}
}

func TestNewConfig(t *testing.T) {
	customTimeout := 15 * time.Second
	customAppId := "custom_app"

	cfg := NewConfig(
		WithRequestTimeout(customTimeout),
		WithAppId(customAppId),
	)

	if cfg.RequestTimeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, cfg.RequestTimeout)
	}

	if cfg.RequestConfig.AppId != customAppId {
		t.Errorf("Expected AppId %s, got %s", customAppId, cfg.RequestConfig.AppId)
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := NewDefaultConfig()

	// Test WithRequestTimeout
	newTimeout := 20 * time.Second
	opt := WithRequestTimeout(newTimeout)
	opt(cfg)
	if cfg.RequestTimeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, cfg.RequestTimeout)
	}

	// Test WithAppId
	newAppId := "test_app_id"
	opt = WithAppId(newAppId)
	opt(cfg)
	if cfg.RequestConfig.AppId != newAppId {
		t.Errorf("Expected AppId %s, got %s", newAppId, cfg.RequestConfig.AppId)
	}

	// Test WithCurBuvid
	newBuvid := "test_buvid"
	opt = WithCurBuvid(newBuvid)
	opt(cfg)
	if cfg.RequestConfig.CurBuvid != newBuvid {
		t.Errorf("Expected CurBuvid %s, got %s", newBuvid, cfg.RequestConfig.CurBuvid)
	}
}

func TestConstants(t *testing.T) {
	if AppkeyPcr == "" {
		t.Error("AppkeyPcr should not be empty")
	}

	if DefaultCaptchaType != "1" {
		t.Errorf("Expected DefaultCaptchaType '1', got %s", DefaultCaptchaType)
	}

	if DefaultCaptchaVersion != "1" {
		t.Errorf("Expected DefaultCaptchaVersion '1', got %s", DefaultCaptchaVersion)
	}

	if MethodPost != "POST" {
		t.Errorf("Expected MethodPost 'POST', got %s", MethodPost)
	}
}
