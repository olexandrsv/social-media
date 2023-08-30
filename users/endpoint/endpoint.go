package endpoint

import (
	"context"
	"social-media/users"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser        endpoint.Endpoint
	Login             endpoint.Endpoint
	GetUser           endpoint.Endpoint
	UpdateUser        endpoint.Endpoint
	GetLoginsByInfo   endpoint.Endpoint
	FollowUser        endpoint.Endpoint
	GetFollowedLogins endpoint.Endpoint
}

func NewEndpoints(s users.Service) Endpoints {
	return Endpoints{
		CreateUser:        makeCreateUserEndpoint(s),
		Login:             makeLoginEndpoint(s),
		GetUser:           makeGetUserEndpoint(s),
		UpdateUser:        makeUpdateUserEndpoint(s),
		GetLoginsByInfo:   makeLoginsByInfoEndpoint(s),
		FollowUser:        makeFollowUserEndpoint(s),
		GetFollowedLogins: makeGetFollowedLogins(s),
	}
}

func makeCreateUserEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserReq)
		token, err := s.CreateUser(req.Login, req.FirstName, req.SecondName, req.Password)
		if err != nil {
			return AuthResp{token, err.Error()}, err
		}
		return AuthResp{token, ""}, nil
	}
}

func makeLoginEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginReq)
		token, err := s.Login(req.Login, req.Password)
		if err != nil {
			return AuthResp{token, err.Error()}, nil
		}
		return AuthResp{token, ""}, nil
	}
}

func makeGetUserEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserReq)
		user, err := s.GetUser(req.Login)
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
}

func makeUpdateUserEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserReq)
		err := s.UpdateUser(req.Token, req.FirstName, req.SecondName, req.Bio, req.Interests)
		if err != nil {
			return UpdateUserResp{Error: err.Error()}, nil
		}
		return UpdateUserResp{Error: ""}, nil
	}
}

func makeLoginsByInfoEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserByInfoReq)
		logins, err := s.GetLoginsByInfo(req.Info)
		if err != nil {
			return nil, err
		}
		return LoginsResp{logins}, nil
	}
}

func makeFollowUserEndpoint(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FollowUserReq)
		err := s.FollowUser(req.Token, req.Login)
		return nil, err
	}
}

func makeGetFollowedLogins(s users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenReq)
		logins, err := s.GetFollowedLogins(req.Token)
		if err != nil{
			return nil, err
		}
		return LoginsResp{logins}, nil
	}
}
