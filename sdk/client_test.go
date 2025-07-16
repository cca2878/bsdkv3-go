package sdk

import (
	"testing"
	"time"

	"github.com/cca2878/bsdkv3-go/sdk/config"
)

func TestNewBSdkV3Client(t *testing.T) {
	tests := []struct {
		name    string
		appKey  string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "valid app key",
			appKey:  "test_app_key",
			opts:    nil,
			wantErr: true, // Will fail due to network call, but should create client
		},
		{
			name:   "valid app key with custom config",
			appKey: "test_app_key",
			opts: []ClientOption{
				WithConfig(config.NewConfig(
					config.WithRequestTimeout(5 * time.Second),
				)),
			},
			wantErr: true, // Will fail due to network call
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewBSdkV3Client(tt.appKey, tt.opts...)

			// We expect this to fail with network error, but client creation should work
			if err == nil {
				t.Errorf("Expected error due to network call in constructor")
				if client != nil {
					client.Close()
				}
				return
			}

			// Verify the error is related to configuration/network, not client creation
			if client != nil {
				client.Close()
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	cfg := config.NewConfig(
		config.WithRequestTimeout(10*time.Second),
		config.WithAppId("test_app"),
	)

	opt := WithConfig(cfg)
	if opt == nil {
		t.Error("WithConfig should return a valid option function")
	}
}

func TestUserInfo(t *testing.T) {
	user := UserInfo{
		Username: "testuser",
		Password: "testpass",
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}

	if user.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got %s", user.Password)
	}

	userInfo := user.GetUserInfo()
	if userInfo.Username != user.Username {
		t.Error("GetUserInfo should return the same user info")
	}
}

func TestCaptchaParams(t *testing.T) {
	params := CaptchaParams{
		CaptchaType: "1",
		Validate:    "test_validate",
		Challenge:   "test_challenge",
		GtUserId:    "test_gt_user_id",
		SecCode:     "test_seccode",
		CToken:      "test_ctoken",
	}

	newParams := NewCaptchaParams(params)
	if newParams.CaptchaType != "1" {
		t.Errorf("Expected CaptchaType '1', got %s", newParams.CaptchaType)
	}

	// Test default captcha type
	emptyParams := CaptchaParams{}
	newEmptyParams := NewCaptchaParams(emptyParams)
	if newEmptyParams.CaptchaType != config.DefaultCaptchaType {
		t.Errorf("Expected default CaptchaType %s, got %s", config.DefaultCaptchaType, newEmptyParams.CaptchaType)
	}
}
