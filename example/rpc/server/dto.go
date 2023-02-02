package main

type UserService struct {
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
}
