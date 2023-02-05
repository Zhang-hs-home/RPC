package RPC

import (
	"RPC/message"
	"RPC/serialize"
	"RPC/serialize/json"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
)

type Server struct {
	services    map[string]reflectionStub
	serializers []serialize.Serializer
}

func NewServer() *Server {
	res := &Server{
		services: map[string]reflectionStub{}, // 存储服务端所提供的服务
		// 一个byte可表示的最大个数为256
		serializers: make([]serialize.Serializer, 256), // 存储服务端支持的序列化协议
	}
	res.RegisterSerializer(json.SerializerJson{}) // 默认支持json序列化

	return res
}

func (s *Server) MustRegister(service Service) {
	if err := s.RegisterService(service); err != nil {
		panic(err)
	}
}

func (s *Server) RegisterService(service Service) error {
	s.services[service.Name()] = reflectionStub{
		value:       reflect.ValueOf(service),
		serializers: s.serializers,
	}
	return nil
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Code()] = serializer
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connection failed, err: ", err)
			continue
		}
		log.Println("receive connection: ", conn.RemoteAddr(), " -> ", conn.LocalAddr())
		go func() {
			if er := s.HandleConn(conn); er != nil {
				log.Println("handle connection failed, err: ", er)
				conn.Close()
				return
			}
		}()
	}

}

func (s *Server) HandleConn(conn net.Conn) error {
	for {
		// 读请求
		// 执行
		// 写回响应

		reqMsg, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		req := message.DecodeReq(reqMsg)
		log.Println(req)
		resp := &message.Response{
			MessageId:  req.MessageId,
			Version:    req.Version,
			Compressor: req.Compressor,
			Serializer: req.Serializer,
		}
		// ping探活请求处理
		if req.Ping == PingPong {
			resp.Pong = PingPong
			resp.CalcHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				return err
			}

			continue
		}

		// 找到本地对应的服务
		service, ok := s.services[req.ServiceName]
		if !ok {
			resp.Error = []byte("找不到对应服务")
			resp.CalcHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				return err
			}

			continue
		}

		// 找到方法 这里只能通过反射 没别的路子

		// 把参数传进去
		ctx := context.Background()
		// 需要把req.Arg 赋值给methodReq Arg是map[string]interface{}类型
		// 编码resp返回

		data, err := service.invoke(ctx, req)
		if err != nil {
			resp.Error = []byte(err.Error())
			resp.CalcHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				return err
			}

			continue
		}

		resp.Data = data
		resp.BodyLength = uint32(len(data))
		resp.CalcHeadLength()
		bitFlow := message.EncodeResp(resp)
		_, err = conn.Write(bitFlow)
		if err != nil {
			return err
		}

	}
}

type reflectionStub struct {
	value       reflect.Value
	serializers []serialize.Serializer
}

// 反射相关的封装在这里,这个方法的作用是：调用本地的服务，拿到返回值，序列化后返回，即response body就是他
func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	methodName := req.MethodName
	data := req.Data

	serializer := s.serializers[req.Serializer]
	if serializer == nil {
		return nil, errors.New("不支持的序列化协议")
	}

	method := s.value.MethodByName(methodName)
	if !method.IsValid() {
		return nil, errors.New(fmt.Sprintf("%s%s", "服务下不存在该方法，方法名为：", methodName))
	}

	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := serializer.Decode(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}

	return serializer.Encode(res[0].Interface())
}
