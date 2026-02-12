package server

import (
	"context"
	"net"
	"social-media/api/pb/auth"
	"social-media/internal/authentication/service"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	"google.golang.org/grpc"
)

type server struct {
	auth.UnimplementedAuthenticateServer
	service service.Service
	srv     *grpc.Server
}

func newServer(service service.Service, s *grpc.Server) *server {
	return &server{
		service: service,
		srv:     s,
	}
}

func NewGRPCServer(service service.Service) *server {
	s := grpc.NewServer()
	srv := newServer(service, s)
	auth.RegisterAuthenticateServer(s, srv)
	return srv
}

func (s *server) GenerateJWT(ctx context.Context, req *auth.GenerateJWTReq) (*auth.GenerateJWTResp, error) {
	token, err := s.service.GenerateToken(int(req.Id), req.Login)

	return &auth.GenerateJWTResp{
		Token: token,
		Err:   convertToError(err),
	}, nil
}

func (s *server) ValidateJWT(ctx context.Context, req *auth.ValidateJWTReq) (*auth.ValidateJWTResp, error) {
	id, login, err := s.service.ValidateToken(req.Token)

	return &auth.ValidateJWTResp{
		Id:    int64(id),
		Login: login,
		Err:   convertToError(err),
	}, nil
}

func (s *server) GenerateSignedUrl(ctx context.Context, req *auth.GenerateSignedUrlReq) (*auth.GenerateSignedUrlResp, error) {
	token, err := s.service.GenerateSignedUrl(req.FileID)

	return &auth.GenerateSignedUrlResp{
		Token: token,
		Err:   convertToError(err),
	}, nil
}

func (s *server) ValidateSignedUrl(ctx context.Context, req *auth.ValidateSignedUrlReq) (*auth.ValidateSignedUrlResp, error) {
	fileID, err := s.service.ValidateSignedUrl(req.Token)

	return &auth.ValidateSignedUrlResp{
		FileID: fileID,
		Err:    convertToError(err),
	}, nil
}

func convertToError(err error) *auth.Error {
	if err == nil {
		return nil
	}

	customError, ok := err.(common.Error)
	if !ok {
		customError = common.ErrInternal
	}

	return &auth.Error{
		Code:    int64(customError.Code()),
		Message: customError.Message(),
	}
}

func (s *server) Run() {
	listener, err := net.Listen("tcp", ":"+config.App.AuthService.Port)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		log.Error(err)
		panic(err)
	}
}
