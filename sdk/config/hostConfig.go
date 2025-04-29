package config

import (
	"strings"
	"sync"
	"time"
)

type hostType string

const (
	HostTypeInitConf   hostType = "init_conf"   // 初始配置
	HostTypeLoginHttps hostType = "login_https" // 登录 HTTPS API
)

var DefaultHosts = map[hostType]string{
	HostTypeInitConf:   defaultInitConfHost,
	HostTypeLoginHttps: defaultLoginHttpsHost,
}

type HostConfig struct {
	mu           sync.RWMutex
	hosts        map[hostType][]string
	lastUpdate   time.Time
	updatePeriod time.Duration
}

// 使用包级变量和init函数初始化单例
var (
	hostConfig *HostConfig
)

func init() {
	hostConfig = &HostConfig{
		hosts:        make(map[hostType][]string),
		updatePeriod: 5 * time.Minute,
	}

	// 初始化默认值
	for hostType, defaultHost := range DefaultHosts {
		hostConfig.hosts[hostType] = []string{defaultHost}
	}
}

func ParseHostsStr(hostsType hostType, hostsStr string) map[hostType][]string {
	return map[hostType][]string{
		hostsType: strings.Split(hostsStr, ","),
	}
}

// GetHostConfig 返回共享的HostConfig实例
func GetHostConfig() *HostConfig {
	return hostConfig
}

func (h *HostConfig) GetHost(hostType hostType) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if hosts, ok := h.hosts[hostType]; ok && len(hosts) > 0 {
		return hosts[0]
	}
	return DefaultHosts[hostType]
}

func (h *HostConfig) UpdateHosts(newHosts map[hostType][]string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for hostType, hosts := range newHosts {
		if len(hosts) > 0 {
			h.hosts[hostType] = hosts
		}
	}
	h.lastUpdate = time.Now()
}

func (h *HostConfig) NeedsUpdate() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return time.Since(h.lastUpdate) > h.updatePeriod
}
