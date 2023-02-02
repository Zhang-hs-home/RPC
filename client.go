package RPC

import (
	"RPC/message"
	"RPC/serialize"
	"context"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"sync/atomic"
	"time"
)

var messageId uint32 = 0

type Service interface {
	Name() string // 用name来寻找服务名称，至于你结构体，结构体里的方法字段叫什么已经解耦了
}

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
		}, // 这个pool缓存的内容是连接
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
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
	defer func() {
		c.coonPool.Put(obj)
	}()
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

func (c *Client) Close() {
	c.coonPool.Release()
}

func (c *Client) InitService(service Service) error { // 校验入参，调用proxy，接收返回值
	// 可以校验，确保它是一个指向结构体的指针

	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		if fieldType.Type.Kind() != reflect.Func {
			continue
		}

		// 替换实现
		fn := reflect.MakeFunc(fieldType.Type, func(args []reflect.Value) (results []reflect.Value) {
			// 你可以对args和results进行校验
			ctx, ok := args[0].Interface().(context.Context)
			if !ok {
				panic("system internal error")
			}
			arg := args[1].Interface()
			// 第一个返回值，是真的返回值，指向GetByIdResp
			outType := fieldType.Type.Out(0)

			// 客户端和服务端的结构体名字可能不一样，但是结构体里的方法字段的名称一般是一样的。如果不一样，也可以在方法后加tag的方式（类似gorm指定列名）
			// 但是一般保证结构体里的方法字段名称一样即可。
			bs, err := c.serializer.Encode(arg)
			if err != nil {
				results = append(results, reflect.Zero(outType))
				results = append(results, reflect.ValueOf(err))
				return
			}
			msgId := atomic.AddUint32(&messageId, 1)
			req := &message.Request{
				Compressor:  0,
				Serializer:  c.serializer.Code(),
				MessageId:   msgId,
				Version:     0,
				BodyLength:  uint32(len(bs)),
				ServiceName: service.Name(),
				MethodName:  fieldType.Name,
				Data:        bs,
			}
			req.CalcHeadLength()

			// 该发请求了，不希望把TCP操作直接放在这里，需要一个invoke操作来抽象网络相关的操作，外面的方法只涉及到marshal入参，调用invoke，unmarshal出参
			// 这里只需要类似 res, err := xxx.Invoke()即可拿到服务端的返回值
			res, err := c.Invoke(ctx, req)

			if err != nil {
				results = append(results, reflect.Zero(outType))
				results = append(results, reflect.ValueOf(err))
				return
			}

			// new出一个结构体
			first := reflect.New(outType.Elem()).Interface()
			// 接下来你需要把数据填进去
			// resp.Data -> first
			if len(res.Data) > 0 {
				err = c.serializer.Decode(res.Data, first)
			}

			results = append(results, reflect.ValueOf(first))
			// 第二个返回值，是error
			if err != nil {
				results = append(results, reflect.ValueOf(err))
			} else {
				results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem())) // 这里是为了给nil带上类型信息，否则反射没法解析
			}
			return
		})
		fieldValue.Set(fn)
	}
	return nil
}
