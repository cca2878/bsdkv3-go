package gateway

import (
	"net/http"
	"testing"
)

// TestInjectedTransportIsUsed 验证 WithHTTPGatewayTransport 注入的 *http.Transport
// 确实被默认网关的 http.Client 采用（共享底层 transport 的前提）。
func TestInjectedTransportIsUsed(t *testing.T) {
	rt := &http.Transport{}
	gw := NewHTTPGateway(WithHTTPGatewayTransport(rt)).(*httpGateway)
	if gw.client.Transport != rt {
		t.Fatalf("注入的 Transport 未被采用: got %p, want %p", gw.client.Transport, rt)
	}
}

// TestProxyIgnoredWhenTransportInjected 验证注入 transport 后 Proxy 选项被忽略
// （代理应由注入的 transport 自带）。
func TestProxyIgnoredWhenTransportInjected(t *testing.T) {
	rt := &http.Transport{}
	gw := NewHTTPGateway(
		WithHTTPGatewayProxy("http://127.0.0.1:8888"),
		WithHTTPGatewayTransport(rt),
	).(*httpGateway)
	if gw.client.Transport != rt {
		t.Fatalf("注入 transport 应优先于 Proxy 选项: got %p, want %p", gw.client.Transport, rt)
	}
}
