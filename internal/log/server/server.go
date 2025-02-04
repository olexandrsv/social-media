package server

import (
	"context"
	"net"
	"social-media/api/pb/log"
	"social-media/internal/common/app/config"
	"social-media/internal/log/service"

	"google.golang.org/grpc"
)

type server struct {
	srv *grpc.Server
	log.UnimplementedLogServer
	service service.Service
}

func NewGRPCServer(service service.Service) *server {
	s := grpc.NewServer()
	srv := &server{
		srv:     s,
		service: service,
	}
	log.RegisterLogServer(s, srv)
	return srv
}

func (s *server) Error(ctx context.Context, req *log.LogRequest) (*log.Empty, error) {
	s.service.Error(req.Msg)
	return &log.Empty{}, nil
}

func (s *server) Info(ctx context.Context, req *log.LogRequest) (*log.Empty, error) {
	s.service.Info(req.Msg)
	return &log.Empty{}, nil
}

func (s *server) Run() {
	listener, err := net.Listen("tcp", ":"+config.App.LogService.Port)
	if err != nil {
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		panic(err)
	}
}
