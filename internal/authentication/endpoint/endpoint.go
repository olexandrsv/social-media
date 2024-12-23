package endpoint

import (
	"context"
	"social-media/api/pb/auth"
	"social-media/internal/authentication/service"
)

type Endpoints struct {
	auth.UnimplementedAuthenticateServer
	service service.Service
}

func NewEndpoints(s service.Service) Endpoints {
	return Endpoints{service: s}
}

func (e Endpoints) GenerateJWT(ctx context.Context, req *auth.GenerateJWTReq) (*auth.GenerateJWTResp, error) {
	token, err := e.service.GenerateToken(int(req.Id), req.Login)
	if err != nil {
		return nil, err
	}
	return &auth.GenerateJWTResp{Token: token}, nil
}

func (e Endpoints) ValidateJWT(ctx context.Context, req *auth.ValidateJWTReq) (*auth.ValidateJWTResp, error) {
	id, login, err := e.service.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &auth.ValidateJWTResp{
		Id:    int64(id),
		Login: login,
	}, nil
}
