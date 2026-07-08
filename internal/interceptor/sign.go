package interceptor

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/url"
	"sort"

	"github.com/cca2878/bsdkv3-go/transport"
)

func calcSign(appKey string, values url.Values) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := md5.New()
	for _, k := range keys {
		_, _ = io.WriteString(h, values.Get(k))
	}
	_, _ = io.WriteString(h, appKey)

	return hex.EncodeToString(h.Sum(nil))
}

func NewSignInterceptor(appKey string) Interceptor {
	return func(ctx context.Context, req *transport.Request, next Invoker) (*transport.Response, error) {
		// 清除上一轮重试可能残留的签名，避免旧 sign 被计入本轮哈希导致签名错误。
		req.FormParams.Del("sign")
		req.FormParams.Set("sign", calcSign(appKey, req.FormParams))
		return next(ctx, req)
	}
}
