package interceptor

import (
	"context"

	"github.com/cca2878/bsdkv3-go/config"
	"github.com/cca2878/bsdkv3-go/transport"
	"github.com/go-playground/form/v4"
)

func NewCommonParamsInterceptor(params config.BaseReqParams) Interceptor {
	cachedValues, err := form.NewEncoder().Encode(params)
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, req *transport.Request, next Invoker) (*transport.Response, error) {
		for key, values := range cachedValues {
			if _, ok := req.FormParams[key]; !ok {
				newSlice := make([]string, len(values))
				copy(newSlice, values)
				req.FormParams[key] = newSlice
			}
		}
		return next(ctx, req)
	}
}
