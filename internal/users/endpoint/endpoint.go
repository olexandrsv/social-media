package endpoint

import (
	"context"
	"errors"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/users/service"
)

type Endpoints interface {
	CreateUser(ctx context.Context, request interface{}) (interface{}, error)
	Login(ctx context.Context, request interface{}) (interface{}, error)
	GetUser(ctx context.Context, request interface{}) (interface{}, error)
	UpdateUser(ctx context.Context, request interface{}) (interface{}, error)
	GetUsersByInfo(ctx context.Context, request interface{}) (interface{}, error)
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
	if !ok {
		log.Error(errors.New("can't assign to CreateUserReq"))
		return nil, common.ErrInternal
	}

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	v.NotLess("password", req.Password, 8)
	if err := v.Err(); err != nil {
		return nil, err
	}

	token, id, err := e.service.CreateUser(req.Login, req.Name, req.Surname, req.Password)
	if err != nil {
		return nil, err
	}
	return AuthResp{
		ID:    id,
		Token: token,
	}, nil
}

func (e usersEndpoints) Login(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(LoginReq)
	if !ok {
		log.Error(errors.New("can't assign to LoginReq"))
		return nil, common.ErrInternal
	}

	v := common.NewValidator()
	v.NotEmpty("login", req.Login)
	v.NotEmpty("password", req.Password)
	if err := v.Err(); err != nil {
		return nil, err
	}

	token, id, err := e.service.Login(req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	return AuthResp{
		ID:    id,
		Token: token,
	}, nil
}

func (e usersEndpoints) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(GetUserReq)
	if !ok {
		log.Error(errors.New("can't assign to GetUserReq"))
		return nil, common.ErrInternal
	}

	user, err := e.service.GetUser(req.Token, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetUserResp{
		Login:     user.Login(),
		Name:      user.Name(),
		Surname:   user.Surname(),
		Bio:       user.Bio(),
		Interests: user.Interests(),
	}, nil
}

func (e usersEndpoints) UpdateUser(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(UpdateUserReq)
	if !ok {
		log.Error(errors.New("can't assign to UpdateUserReq"))
		return nil, common.ErrInternal
	}
	err := e.service.UpdateUser(req.Token, req.Name, req.Surname, req.Bio, req.Interests)
	if err != nil {
		return UpdateUserResp{Error: err.Error()}, nil
	}
	return UpdateUserResp{Error: ""}, nil
}

func (e usersEndpoints) GetUsersByInfo(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(GetLoginsByInfoReq)
	if !ok {
		log.Error(errors.New("can't assign to GetLoginsByInfoReq"))
		return nil, common.ErrInternal
	}
	users, err := e.service.GetUsersByInfo(req.Info)
	if err != nil {
		return nil, err
	}
	userModels := make([]UserModel, 0, len(users))
	for _, user := range users {
		userModels = append(userModels, UserModel{
			ID:    user.ID(),
			Login: user.Login(),
		})
	}
	return userModels, nil
}

func (e usersEndpoints) FollowUser(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(FollowUserReq)
	if !ok {
		log.Error(errors.New("can't assign to FollowUserReq"))
		return nil, common.ErrInternal
	}
	err := e.service.FollowUser(req.Token, req.Login)
	return nil, err
}

func (e usersEndpoints) GetFollowedLogins(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(TokenReq)
	if !ok {
		log.Error(errors.New("can't assign to TokenReq"))
		return nil, common.ErrInternal
	}
	logins, err := e.service.GetFollowedLogins(req.Token)
	if err != nil {
		return nil, err
	}
	return LoginsResp{logins}, nil
}
