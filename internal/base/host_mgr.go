package base

import "sync"

// host_manager.go
type HostMgr struct {
	mu      sync.RWMutex
	hosts   []string
	currIdx int // 记录当前正在使用的 Host 索引
}

// NewHostManager 初始化
func NewHostManager(initialHosts []string) *HostMgr {
	return &HostMgr{
		hosts:   initialHosts,
		currIdx: 0,
	}
}

func (m *HostMgr) GetTotalHosts() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.hosts)
}

// GetCurrentHost 获取当前应该使用的 Host
func (m *HostMgr) GetCurrentHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.hosts) == 0 {
		return ""
	}
	return m.hosts[m.currIdx]
}

// MarkFailed 报告某个 Host 挂了，切换到下一个
func (m *HostMgr) MarkFailed(failedHost string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.hosts) == 0 {
		return
	}
	// 核心防御逻辑：只有当报告失败的 Host 等于当前 Host 时才切换。
	// 这是为了防止并发请求时，A和B同时失败，导致 currIdx 被加了两次，跳过了一个原本可用的 Host。
	currentHost := m.hosts[m.currIdx]
	if currentHost == failedHost {
		m.currIdx = (m.currIdx + 1) % len(m.hosts)
	}
}

// UpdateHosts 更新服务端下发的新列表
func (m *HostMgr) UpdateHosts(newHosts []string) {
	if len(newHosts) == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hosts = newHosts
	m.currIdx = 0 // 拿到新列表后，重置回第一个开始用
}
