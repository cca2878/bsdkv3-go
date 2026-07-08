package interceptor

import (
	"context"
	"fmt"

	"github.com/cca2878/bsdkv3-go/config"
	"github.com/cca2878/bsdkv3-go/transport"
	"github.com/go-playground/form/v4"
)

func NewCommonParamsInterceptor(params config.BaseReqParams) Interceptor {
	// 公共参数是一组扁平字符串字段，正常绝不会编码失败；一旦失败属于不变量被破坏，
	// panic 附带上下文，方便集成方定位是「公共参数编码」环节而非底层裸错误。
	cachedValues, err := form.NewEncoder().Encode(params)
	if err != nil {
		panic(fmt.Errorf("bsdkv3: 编码公共参数(BaseReqParams)失败: %w", err))
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
