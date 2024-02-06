package endpoint

import (
	"context"
	"social-media/common"
	"social-media/users"
)

type Endpoints struct {
	service users.Service
}

func NewEndpoints(s users.Service) Endpoints {
	return Endpoints{
		service: s,
	}
}

func (e Endpoints) CreateUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(CreateUserReq)

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	v.NotLess("password", req.Password, 8)
	if err := v.Err(); err != nil {
		return nil, err
	}

	token, err := e.service.CreateUser(req.Login, req.FirstName, req.SecondName, req.Password)
	if err != nil {
		return nil, err
	}
	return AuthResp{token}, nil
}

func (e Endpoints) Login(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(LoginReq)

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	v.NotEmpty("password", req.Password)
	if err := v.Err(); err != nil {
		return nil, err
	}

	token, err := e.service.Login(req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	return AuthResp{token}, nil
}

func (e Endpoints) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(GetUserReq)

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	if err := v.Err(); err != nil {
		return nil, err
	}

	user, err := e.service.GetUser(req.Login)
	if err != nil {
		return nil, err
	}
	return &GetUserResp{
		Login:      user.Login,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Bio:        user.Bio,
		Interests:  user.Interests,
	}, nil
}

func (e Endpoints) UpdateUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(UpdateUserReq)
	err := e.service.UpdateUser(req.Token, req.FirstName, req.SecondName, req.Bio, req.Interests)
	if err != nil {
		return UpdateUserResp{Error: err.Error()}, nil
	}
	return UpdateUserResp{Error: ""}, nil
}

func (e Endpoints) GetLoginsByInfo(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(GetUserByInfoReq)
	logins, err := e.service.GetLoginsByInfo(req.Info)
	if err != nil {
		return nil, err
	}
	return LoginsResp{logins}, nil
}

func (e Endpoints) FollowUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(FollowUserReq)
	err := e.service.FollowUser(req.Token, req.Login)
	return nil, err
}

func (e Endpoints) GetFollowedLogins(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(TokenReq)
	logins, err := e.service.GetFollowedLogins(req.Token)
	if err != nil {
		return nil, err
	}
	return LoginsResp{logins}, nil
}

func (e Endpoints) Log(err error){
	e.service.Log(err)
}