package main

import (
	"RPC"
	"RPC/serialize/json"
	"RPC/serialize/protobuf"
)

func main() {
	// 启动server示例
	srv := RPC.NewServer()
	// 注册server服务
	srv.MustRegister(&UserService{})
	srv.MustRegister(&UserParentService{})
	// 注册server支持的序列化协议
	srv.RegisterSerializer(json.SerializerJson{})
	srv.RegisterSerializer(protobuf.SerializerProto{})
	if err := srv.Start("localhost:8080"); err != nil {
		panic(err)
	}

}
