package config

import (
	"testing"
	"time"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg == nil {
		t.Fatal("NewDefaultConfig() returned nil")
	}

	// Test some expected default values
	if cfg.RequestTimeout != DefaultRequestTimeout*time.Second {
		t.Errorf("Expected timeout %v, got %v", DefaultRequestTimeout*time.Second, cfg.RequestTimeout)
	}

	if cfg.RequestConfig.SdkType == "" {
		t.Error("Expected non-empty SdkType")
	}

	if cfg.RequestConfig.Platform == "" {
		t.Error("Expected non-empty Platform")
	}
}

func TestNewConfig(t *testing.T) {
	customTimeout := 10 * time.Second

	cfg := NewConfig(
		WithRequestTimeout(customTimeout),
		WithAppId("test_app"),
		WithGameId("test_game"),
	)

	if cfg.RequestTimeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, cfg.RequestTimeout)
	}

	if cfg.RequestConfig.AppId != "test_app" {
		t.Errorf("Expected AppId 'test_app', got '%s'", cfg.RequestConfig.AppId)
	}

	if cfg.RequestConfig.GameId != "test_game" {
		t.Errorf("Expected GameId 'test_game', got '%s'", cfg.RequestConfig.GameId)
	}
}

func TestWithConfigOptions(t *testing.T) {
	tests := []struct {
		name    string
		option  Option
		checkFn func(*Config) bool
	}{
		{
			name:   "WithRequestTimeout",
			option: WithRequestTimeout(5 * time.Second),
			checkFn: func(c *Config) bool {
				return c.RequestTimeout == 5*time.Second
			},
		},
		{
			name:   "WithCurBuvid",
			option: WithCurBuvid("test_buvid"),
			checkFn: func(c *Config) bool {
				return c.RequestConfig.CurBuvid == "test_buvid"
			},
		},
		{
			name:   "WithAppId",
			option: WithAppId("test_app_id"),
			checkFn: func(c *Config) bool {
				return c.RequestConfig.AppId == "test_app_id"
			},
		},
		{
			name:   "WithPlatform",
			option: WithPlatform("test_platform"),
			checkFn: func(c *Config) bool {
				return c.RequestConfig.Platform == "test_platform"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfig(tt.option)
			if !tt.checkFn(cfg) {
				t.Errorf("Option %s did not apply correctly", tt.name)
			}
		})
	}
}
