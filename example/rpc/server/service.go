package main

import "context"

func (u *UserService) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Name: "tom",
	}, nil
}
