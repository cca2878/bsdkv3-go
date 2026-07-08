package interceptor

import (
	"context"
	"maps"
	"strconv"
	"time"

	"github.com/cca2878/bsdkv3-go/transport"
)

func NewStampInterceptor() Interceptor {
	return func(ctx context.Context, req *transport.Request, next Invoker) (*transport.Response, error) {
		domain := req.URL.Hostname()
		req.FormParams.Set("domain", domain)
		req.FormParams.Set("original_domain", domain)
		req.FormParams.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
		maps.Copy(req.Headers, map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   "Mozilla/5.0 BSGameSDK",
			"cversion":     "1",
		})
		return next(ctx, req)
	}
}
