package bsdkv3

import "sync"

// host_manager.go
type hostMgr struct {
	mu      sync.RWMutex
	hosts   []string
	currIdx int // 记录当前正在使用的 Host 索引
}

// newHostManager 初始化
func newHostManager(initialHosts []string) *hostMgr {
	return &hostMgr{
		hosts:   initialHosts,
		currIdx: 0,
	}
}

// getCurrentHost 获取当前应该使用的 Host
func (m *hostMgr) getCurrentHost() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.hosts) == 0 {
		return ""
	}
	// 利用取模运算实现循环读取
	return m.hosts[m.currIdx%len(m.hosts)]
}

// markFailed 报告某个 Host 挂了，切换到下一个
func (m *hostMgr) markFailed(failedHost string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.hosts) == 0 {
		return
	}
	// 核心防御逻辑：只有当报告失败的 Host 等于当前 Host 时才切换。
	// 这是为了防止并发请求时，A和B同时失败，导致 currIdx 被加了两次，跳过了一个原本可用的 Host。
	currentHost := m.hosts[m.currIdx%len(m.hosts)]
	if currentHost == failedHost {
		m.currIdx = (m.currIdx + 1) % len(m.hosts)
	}
}

// updateHosts 更新服务端下发的新列表
func (m *hostMgr) updateHosts(newHosts []string) {
	if len(newHosts) == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hosts = newHosts
	m.currIdx = 0 // 拿到新列表后，重置回第一个开始用
}
