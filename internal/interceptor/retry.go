package interceptor

import (
	"context"
	"fmt"
	"net/url"

	"github.com/cca2878/bsdkv3-go/internal/apierr"
	"github.com/cca2878/bsdkv3-go/internal/base"
	"github.com/cca2878/bsdkv3-go/transport"
)

func NewRetryInterceptor(maxRetries int, hostMgr *base.HostMgr) Interceptor {
	// maxRetries 语义为「每个 host 的尝试次数」，至少 1（0/负数会让下方 for 循环一次不跑，
	// 直接返回 (nil, nil)，进而在 execAPI 里 resp.Body 空指针 panic）。
	if maxRetries < 1 {
		maxRetries = 1
	}
	return func(ctx context.Context, req *transport.Request, next Invoker) (*transport.Response, error) {
		// 保存原始相对 URL：每轮都基于它重新解析，否则首轮解析成绝对 URL 后，
		// 后续 ResolveReference(绝对URL) 会原样返回，导致换 host 不生效。
		relURL := req.URL
		var lastErr error
		triedHostsNum := 0

		for i := 0; i < maxRetries; i++ {
			// 1. 动态决定路由
			currentHost := hostMgr.GetCurrentHost()
			baseURL, err := url.Parse(currentHost)
			if err != nil {
				return nil, err
			}
			// 上层应该只传相对URL
			req.URL = baseURL.ResolveReference(relURL)

			// 2. 向下放行
			resp, err := next(ctx, req)

			// 3. 错误捕获与路由惩罚
			if err != nil || (resp != nil && resp.StatusCode >= 500) {
				if i == maxRetries-1 {
					// 触发切换
					hostMgr.MarkFailed(currentHost)
					triedHostsNum++
					if triedHostsNum >= hostMgr.GetTotalHosts() {
						// 所有主机都尝试过了，放弃重试。
						// 若无传输错误（纯 5xx），也要合成一个 error：否则上层会把 5xx 错误页
						// 当成成功响应去反序列化，掩盖真正的服务端故障。
						if err == nil && resp != nil {
							return resp, fmt.Errorf("%w: 末次 HTTP 状态码 %d", apierr.ErrAllHostsFailed, resp.StatusCode)
						}
						return resp, err
					}
					i = -1 // 重置重试计数器
				}
				lastErr = err
				continue
			}

			// 成功，直接返回
			return resp, nil
		}
		return nil, lastErr
	}
}
