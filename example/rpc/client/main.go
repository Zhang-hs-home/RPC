package main

import (
	"RPC"
	"RPC/serialize/json"
	"context"
	"fmt"
	"log"
)

func main() {
	// 初始化客户端，期望采用json序列化协议
	c, err := RPC.NewClient("localhost:8080", json.SerializerJson{})
	if err != nil {
		log.Println(err)
		return
	}
	// 初始化服务，调用服务端的"user-service"服务 下的GetById方法
	us := &UserService{}
	err = c.InitService(us)
	if err != nil {
		log.Println(err)
		return
	}
	res, err := us.GetById(context.Background(), &GetByIdReq{
		Id: 123,
	})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("%+v", res))

}
