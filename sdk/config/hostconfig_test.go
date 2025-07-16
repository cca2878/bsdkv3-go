package config

import (
	"testing"
)

func TestGetHostConfig(t *testing.T) {
	hostConfig := GetHostConfig()
	if hostConfig == nil {
		t.Fatal("GetHostConfig() returned nil")
	}
}

func TestHostConfig_GetHost(t *testing.T) {
	hostConfig := GetHostConfig()

	tests := []struct {
		name     string
		hostType HostType
		wantHost string
	}{
		{
			name:     "init conf host",
			hostType: HostTypeInitConf,
			wantHost: defaultInitConfHost,
		},
		{
			name:     "login https host",
			hostType: HostTypeLoginHttps,
			wantHost: defaultLoginHttpsHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hostConfig.GetHost(tt.hostType)
			if got != tt.wantHost {
				t.Errorf("GetHost() = %v, want %v", got, tt.wantHost)
			}
		})
	}
}

func TestHostConfig_UpdateHosts(t *testing.T) {
	hostConfig := GetHostConfig()

	// Test updating hosts
	newHosts := map[HostType][]string{
		HostTypeLoginHttps: {"https://new-host1.example.com", "https://new-host2.example.com"},
	}

	hostConfig.UpdateHosts(newHosts)

	// Check that the first host is now the new one
	got := hostConfig.GetHost(HostTypeLoginHttps)
	expected := "https://new-host1.example.com"
	if got != expected {
		t.Errorf("After UpdateHosts(), GetHost() = %v, want %v", got, expected)
	}

	// Test with empty hosts (should not update)
	emptyHosts := map[HostType][]string{
		HostTypeLoginHttps: {},
	}

	hostConfig.UpdateHosts(emptyHosts)

	// Should still have the previous host
	got = hostConfig.GetHost(HostTypeLoginHttps)
	if got != expected {
		t.Errorf("After UpdateHosts() with empty slice, GetHost() = %v, want %v", got, expected)
	}
}

func TestParseHostsStr(t *testing.T) {
	tests := []struct {
		name      string
		hostType  HostType
		hostsStr  string
		wantCount int
	}{
		{
			name:      "single host",
			hostType:  HostTypeLoginHttps,
			hostsStr:  "https://host1.example.com",
			wantCount: 1,
		},
		{
			name:      "multiple hosts",
			hostType:  HostTypeLoginHttps,
			hostsStr:  "https://host1.example.com,https://host2.example.com,https://host3.example.com",
			wantCount: 3,
		},
		{
			name:      "empty string",
			hostType:  HostTypeLoginHttps,
			hostsStr:  "",
			wantCount: 1, // empty string still creates one element
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseHostsStr(tt.hostType, tt.hostsStr)
			hosts, exists := result[tt.hostType]
			if !exists {
				t.Errorf("ParseHostsStr() did not return expected host type")
				return
			}
			if len(hosts) != tt.wantCount {
				t.Errorf("ParseHostsStr() returned %d hosts, want %d", len(hosts), tt.wantCount)
			}
		})
	}
}
