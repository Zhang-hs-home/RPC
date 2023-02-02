package RPC

import (
	"RPC/message"
	"RPC/serialize"
	"context"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"time"
)

//  最后把这里的所有代码移动到client里了，我没动哈

type Client struct {
	coonPool   pool.Pool
	serializer serialize.Serializer
}

func NewClient(addr string, serializer serialize.Serializer) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap:     100,
		MaxIdle:    30,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		}, // 这个pool缓存的东西是连接
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		}, // 在垃圾回收的时候，它会帮你调用close方法关掉连接
		IdleTimeout: time.Minute,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		coonPool:   p,
		serializer: serializer,
	}, nil

}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	obj, err := c.coonPool.Get()
	if err != nil {
		return nil, err
	}

	conn := obj.(net.Conn)
	// 真正把数据发出去
	data := message.EncodeReq(req)
	wLen, err := conn.Write(data)
	if err != nil {
		return nil, err
	}

	if wLen != len(data) { // 几乎遇不到，如果遇到了，有没有更好的处理方法？
		return nil, errors.New("rpc: 未写入全部数据")
	}

	// 该读数据了
	// 我怎么知道该读取多长？相应的，服务端也会有这个问题。答案：你需要先传入一个长度字段，用于描述本次数据有多大

	respMsg, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}

	return message.DecodeResp(respMsg), nil
}
