package main

import "context"

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type UserParentService struct {
	GetParentById func(ctx context.Context, req *GetParentByIdReq) (*GetParentByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

func (u *UserParentService) Name() string {
	return "user-parent-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
}

type GetParentByIdReq struct {
	Id int
}

type GetParentByIdResp struct {
	Father string `json:"father"`
	Mother string `json:"mother"`
}
