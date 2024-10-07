package endpoint

import (
	"context"
	"social-media/api/pb/auth"
	"social-media/internal/authentication/service"
	"social-media/internal/common"
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
		return nil, common.ErrInternal
	}
	return &auth.GenerateJWTResp{Token: token}, nil
}

func (e Endpoints) ValidateJWT(ctx context.Context, req *auth.ValidateJWTReq) (*auth.ValidateJWTResp, error) {
	id, login, err := e.service.ValidateToken(req.Token)
	if err != nil {
		return nil, common.ErrInvalidToken
	}
	return &auth.ValidateJWTResp{
		Id:    int64(id),
		Login: login,
	}, nil
}

func (e Endpoints) Log(err error) {
	e.service.Log(err)
}
