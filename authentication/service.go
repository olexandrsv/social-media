package authentication

import (
	"context"
	"social-media/api/pb"
)

type server struct{
	pb.UnimplementedAuthenticateServer
}

func NewServer() pb.AuthenticateServer{
	return &server{}
}

func (s *server) GenerateJWT(ctx context.Context, req *pb.GenerateJWTReq) (*pb.GenerateJWTResp, error){
	token, err := GenerateToken(int(req.Id), req.Login)
	if err != nil{
		return nil, err
	}
	return &pb.GenerateJWTResp{Token: token}, nil
}

func (s *server) ValidateJWT(ctx context.Context, req *pb.ValidateJWTReq) (*pb.ValidateJWTResp, error){
	id, login, err := ValidateToken(req.Token)
	if err != nil{
		return nil, err
	}
	return &pb.ValidateJWTResp{
		Id: int64(id),
		Login: login,
	}, nil
}