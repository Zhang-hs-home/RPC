# RPC-Framework

基于go实现的简易RPC框架

## 介绍&&功能

- 参考Dubbo1.0版本协议结构，
  采用自定义协议，分为协议头和协议体，手动对协议进行编码和解码。
- 基于TCP进行网络通信。
- 支持轻松的扩展序列化协议作用于协议体，源代码已支持json，protobuf协议。
- 采用连接池管理客户端连接。
- 未引入注册中心。

## 使用方法

### 服务端
```go
// 首先需要定义服务端所提供的服务，以及服务下的方法，Name方法用于描述服务的名称
type UserService struct {
}

func (u *UserService) Name() string {
	return "user-service"
}

func (u *UserService) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Name: fmt.Sprintf("%d-%s", req.Id, "tom"),
	}, nil
}

// 之后就是启动服务端，注册服务，注册支持的序列化协议，开始监听端口
	
// 创建server实例
srv := RPC.NewServer()
// 注册server服务
srv.MustRegister(&UserService{})
// 注册server支持的序列化协议
srv.RegisterSerializer(json.SerializerJson{})
srv.RegisterSerializer(protobuf.SerializerProto{})
if err := srv.Start("localhost:8080"); err != nil {
	panic(err)
}
```
### 客户端
```go
// 首先需要定义客户端期望调用的远程服务，以及服务下的方法，并实现Name方法用于标识服务名称
type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

// Name方法用于寻找服务端对应的服务
func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
}

// 之后就是初始化客户端实例，初始化stub，之后调用方法
c, err := RPC.NewClient("localhost:8080", json.SerializerJson{})
if err != nil {
    return
}

us := &UserService{}
err = c.InitService(us)
if err != nil {
    return
}

res, err := us.GetById(context.Background(), &GetByIdReq{
    Id: 123,
})
if err != nil {
    return
}

fmt.Println(fmt.Sprintf("%+v", res))
c.Close()
```

更加详细的示例请看example文件夹里的内容

## 注意事项

- 我们强制规定，每一个方法的入参必须由两个参数组成，上下文+request结构体，
返回值也必须由两个参数组成，response结构体+error。
  
  
欢迎各位fork，star，或者提出宝贵的意见！
  

  
