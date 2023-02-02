package main

type UserService struct {
}

func (u *UserService) Name() string {
	return "user-service"
}

type UserParentService struct {
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
