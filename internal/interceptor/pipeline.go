package interceptor

import (
	"context"

	"github.com/cca2878/bsdkv3-go/transport"
)

// Invoker 定义了向下一层传递控制权的动作
type Invoker func(ctx context.Context, req *transport.Request) (*transport.Response, error)

// Interceptor 定义了流水线上的独立处理单元
type Interceptor func(ctx context.Context, req *transport.Request, next Invoker) (*transport.Response, error)

type Pipeline struct {
	gateway      transport.Gateway
	interceptors []Interceptor
	chain        Invoker
}

func NewPipeline(gw transport.Gateway) *Pipeline {
	return &Pipeline{
		gateway:      gw,
		interceptors: make([]Interceptor, 0),
		chain:        gw.Do,
	}
}

func (p *Pipeline) composeChain() {
	p.chain = p.gateway.Do
	for i := len(p.interceptors) - 1; i >= 0; i-- {
		currentInterceptor := p.interceptors[i]
		prevNext := p.chain
		p.chain = func(ctx context.Context, req *transport.Request) (*transport.Response, error) {
			return currentInterceptor(ctx, req, prevNext)
		}
	}
}

func (p *Pipeline) Use(interceptor Interceptor) *Pipeline {
	p.interceptors = append(p.interceptors, interceptor)
	p.composeChain()
	return p
}

func (p *Pipeline) Do(ctx context.Context, req *transport.Request) (*transport.Response, error) {
	return p.chain(ctx, req)
}
