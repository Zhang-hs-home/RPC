package RPC

import (
	"RPC/serialize/json"
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

type UserServiceClient struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserServiceClient) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(":8082", json.SerializerJson{})
	require.NoError(t, err)
	us := &UserServiceClient{}
	err = c.InitStub(us)
	require.NoError(t, err)
	resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 100})
	require.NoError(t, err)
	log.Println(resp)
}
