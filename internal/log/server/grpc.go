package server

import (
	"context"
	"net"
	"social-media/api/pb/log"
	"social-media/internal/common/app/config"
	"social-media/internal/log/service"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type grpcServer struct {
	grpcSrv *grpc.Server
	log.UnimplementedLogServer
	service service.Service
}

func newGrpcServer(service service.Service) *grpcServer {
	s := grpc.NewServer()
	srv := &grpcServer{
		grpcSrv: s,
		service: service,
	}
	log.RegisterLogServer(s, srv)
	return srv
}

func (s *grpcServer) Error(ctx context.Context, req *log.LogRequest) (*log.Empty, error) {
	s.service.Error(req.Msg)
	return &log.Empty{}, nil
}

func (s *grpcServer) Info(ctx context.Context, req *log.LogRequest) (*log.Empty, error) {
	s.service.Info(req.Msg)
	return &log.Empty{}, nil
}

func (s *grpcServer) run() {
	listener, err := net.Listen("tcp", ":"+config.App.LogService.GrpcPort)
	if err != nil {
		s.service.Error(errors.WithStack(err).Error())
		return
	}

	if err := s.grpcSrv.Serve(listener); err != nil {
		s.service.Error(errors.WithStack(err).Error())
		return
	}
}
