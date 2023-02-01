package RPC

import (
	"RPC/message"
	"context"
)

// 我们在客户端整了一个proxy，让客户端以为调用这个proxy就是调用了本地的方法
// proxy是客户端和服务端都有的
// 如果想对proxy进一步封装，可以这么写，中间件都是这个类似的设计
/*
type Filter func(ctx context.Context, req *Request) (*Response, error)

type FilterProxy struct {
    Proxy
    filters []Filter
}

func (f FilterProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
    for _,flt := range f.filters {
        flt(ctx, req)
        .....
    }
}
核心就是写别的结构体，这个结构体一定要实现proxy接口
*/

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}
