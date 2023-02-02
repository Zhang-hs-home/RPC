package main

import (
	"RPC"
	"RPC/serialize/json"
	"context"
	"fmt"
	"log"
	"sync"
)

func main() {
	// 初始化客户端，期望采用json序列化协议和服务端通信
	c, err := RPC.NewClient("localhost:8080", json.SerializerJson{})
	if err != nil {
		log.Println(err)
		return
	}
	// 初始化服务，调用服务端的"user-service"服务 下的GetById方法
	us := &UserService{}
	ups := &UserParentService{}
	err = c.InitService(us)
	if err != nil {
		log.Println(err)
		return
	}

	err = c.InitService(ups)
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(20)
	for j := 0; j < 10; j++ {
		i := j
		go func() {
			res, err := us.GetById(context.Background(), &GetByIdReq{
				Id: i,
			})
			if err != nil {
				log.Println(err)
				return
			}
			// 输出服务端返回结果
			fmt.Println(fmt.Sprintf("%+v", res))
			wg.Done()
		}()

		go func() {
			res2, err := ups.GetParentById(context.Background(), &GetParentByIdReq{
				Id: i,
			})
			if err != nil {
				log.Println(err)
				return
			}
			// 输出服务端返回结果
			fmt.Println(fmt.Sprintf("%+v", res2))
			wg.Done()
		}()
	}
	wg.Wait()

	// 关闭客户端
	c.Close()

}
