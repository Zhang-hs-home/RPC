package RPC

import (
	"RPC/message"
	"RPC/serialize/json"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(":8082", json.SerializerJson{})
	require.NoError(t, err)
	us := &UserServiceClient{}
	err = c.InitService(us)
	require.NoError(t, err)
	resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 100})
	require.NoError(t, err)
	log.Println(resp)
}

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name     string
		service  *UserServiceClient
		p        *mockProxy
		wantReq  *message.Request
		wantResp *GetByIdResp
		wantErr  error
	}{
		{
			name: "user-service",
			p: &mockProxy{
				result: []byte(`{"name": "tom"}`),
			},
			wantReq: &message.Request{
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte(`{"id":122}`),
			},
			wantResp: &GetByIdResp{
				Name: "tom",
			},
			service: &UserServiceClient{},
		},
		{
			name: "proxy error",
			p: &mockProxy{
				err: errors.New("mock error"),
			},
			service: &UserServiceClient{},
			wantErr: errors.New("mock error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service, tc.p)
			require.NoError(t, err)
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 111})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			// 断言p的数据
			assert.Equal(t, tc.wantReq, tc.p.req)
			assert.Equal(t, tc.wantResp, resp)
			print(resp)

		})
	}
}

type mockProxy struct {
	req    interface{}
	err    error
	result []byte
}

func (m *mockProxy) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	m.req = req
	return &message.Response{
		Data: m.result,
	}, m.err
}
