package bsdkv3

import (
	"testing"
)

func TestValidateUserInfo(t *testing.T) {
	tests := []struct {
		name      string
		user      UserInfo
		wantError bool
	}{
		{
			name: "valid user info",
			user: UserInfo{
				Username: "testuser",
				Password: "testpass",
				Platform: "android",
				Channel:  "bilibili",
			},
			wantError: false,
		},
		{
			name: "empty username",
			user: UserInfo{
				Username: "",
				Password: "testpass",
			},
			wantError: true,
		},
		{
			name: "whitespace only username",
			user: UserInfo{
				Username: "   ",
				Password: "testpass",
			},
			wantError: true,
		},
		{
			name: "empty password",
			user: UserInfo{
				Username: "testuser",
				Password: "",
			},
			wantError: true,
		},
		{
			name: "whitespace only password",
			user: UserInfo{
				Username: "testuser",
				Password: "   ",
			},
			wantError: true,
		},
		{
			name: "username too long",
			user: UserInfo{
				Username: string(make([]byte, 101)), // 101 characters
				Password: "testpass",
			},
			wantError: true,
		},
		{
			name: "password too long",
			user: UserInfo{
				Username: "testuser",
				Password: string(make([]byte, 201)), // 201 characters
			},
			wantError: true,
		},
		{
			name: "whitespace only platform",
			user: UserInfo{
				Username: "testuser",
				Password: "testpass",
				Platform: "   ",
			},
			wantError: true,
		},
		{
			name: "whitespace only channel",
			user: UserInfo{
				Username: "testuser",
				Password: "testpass",
				Channel:  "   ",
			},
			wantError: true,
		},
		{
			name: "optional platform and channel empty",
			user: UserInfo{
				Username: "testuser",
				Password: "testpass",
				Platform: "",
				Channel:  "",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserInfo(tt.user)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateUserInfo() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAppKey(t *testing.T) {
	tests := []struct {
		name      string
		appKey    string
		wantError bool
	}{
		{
			name:      "valid app key",
			appKey:    "fe8aac4e02f845b8ad67c427d48bfaf1",
			wantError: false,
		},
		{
			name:      "empty app key",
			appKey:    "",
			wantError: true,
		},
		{
			name:      "whitespace only app key",
			appKey:    "   ",
			wantError: true,
		},
		{
			name:      "app key too short",
			appKey:    "fe8aac4e02f845b8ad67c427d48bfaf",
			wantError: true,
		},
		{
			name:      "app key too long",
			appKey:    "fe8aac4e02f845b8ad67c427d48bfaf12",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppKey(tt.appKey)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAppKey() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
