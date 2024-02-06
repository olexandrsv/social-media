package endpoint

import (
	"context"
	"social-media/api/pb"
	"social-media/authentication"
	"social-media/common"
)

type Endpoints struct {
	pb.UnimplementedAuthenticateServer
	service authentication.Service
}

func NewEndpoints(s authentication.Service) Endpoints {
	return Endpoints{service: s}
}

func (e Endpoints) GenerateJWT(ctx context.Context, req *pb.GenerateJWTReq) (*pb.GenerateJWTResp, error) {
	token, err := e.service.GenerateToken(int(req.Id), req.Login)
	if err != nil {
		return nil, common.ErrInternal
	}
	return &pb.GenerateJWTResp{Token: token}, nil
}

func (e Endpoints) ValidateJWT(ctx context.Context, req *pb.ValidateJWTReq) (*pb.ValidateJWTResp, error) {
	id, login, err := e.service.ValidateToken(req.Token)
	if err != nil {
		return nil, common.ErrInvalidToken
	}
	return &pb.ValidateJWTResp{
		Id:    int64(id),
		Login: login,
	}, nil
}

func (e Endpoints) Log(err error) {
	e.service.Log(err)
}
