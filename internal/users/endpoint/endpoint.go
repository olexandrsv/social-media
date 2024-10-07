package endpoint

import (
	"context"
	"errors"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/users/service"
)

type Endpoints interface{
	CreateUser(ctx context.Context, request interface{}) (interface{}, error)
	Login(ctx context.Context, request interface{}) (interface{}, error)
	GetUser(ctx context.Context, request interface{}) (interface{}, error)
	UpdateUser(ctx context.Context, request interface{}) (interface{}, error)
	GetLoginsByInfo(ctx context.Context, request interface{}) (interface{}, error)
	FollowUser(ctx context.Context, request interface{}) (interface{}, error)
	GetFollowedLogins(ctx context.Context, request interface{}) (interface{}, error)
}

type usersEndpoints struct {
	service service.Service
}

func NewEndpoints(s service.Service) Endpoints {
	return usersEndpoints{
		service: s,
	}
}

func (e usersEndpoints) CreateUser(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CreateUserReq)
	if !ok{
		log.Error(errors.New("can't assign to CreateUserReq"))
		return nil, common.ErrInternal
	}

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	v.NotLess("password", req.Password, 8)
	if err := v.Err(); err != nil {
		return nil, err
	}

	token, err := e.service.CreateUser(req.Login, req.Name, req.Surname, req.Password)
	if err != nil {
		return nil, err
	}
	return AuthResp{token}, nil
}

func (e usersEndpoints) Login(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(LoginReq)
	if !ok{
		log.Error(errors.New("can't assign to LoginReq"))
		return nil, common.ErrInternal
	}

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

func (e usersEndpoints) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
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
		Login:      user.Login(),
		Name:  user.Name(),
		Surname: user.Surname(),
		Bio:        user.Bio(),
		Interests:  user.Interests(),
	}, nil
}

func (e usersEndpoints) UpdateUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(UpdateUserReq)
	err := e.service.UpdateUser(req.Token, req.Name, req.Surname, req.Bio, req.Interests)
	if err != nil {
		return UpdateUserResp{Error: err.Error()}, nil
	}
	return UpdateUserResp{Error: ""}, nil
}

func (e usersEndpoints) GetLoginsByInfo(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(GetLoginsByInfoReq)
	logins, err := e.service.GetLoginsByInfo(req.Info)
	if err != nil {
		return nil, err
	}
	return LoginsResp{logins}, nil
}

func (e usersEndpoints) FollowUser(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(FollowUserReq)
	err := e.service.FollowUser(req.Token, req.Login)
	return nil, err
}

func (e usersEndpoints) GetFollowedLogins(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(TokenReq)
	logins, err := e.service.GetFollowedLogins(req.Token)
	if err != nil {
		return nil, err
	}
	return LoginsResp{logins}, nil
}

