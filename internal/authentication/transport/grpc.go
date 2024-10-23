package transport

import (
	"net"
	"social-media/api/pb/auth"
	"social-media/internal/authentication/endpoint"
	"social-media/internal/common/app/log"

	"google.golang.org/grpc"
)

type server struct {
	endpoints endpoint.Endpoints
	srv       *grpc.Server
}

func newServer(e endpoint.Endpoints, s *grpc.Server) *server {
	return &server{e, s}
}

func NewGRPCServer(endpoints endpoint.Endpoints) *server {
	s := grpc.NewServer()
	auth.RegisterAuthenticateServer(s, endpoints)
	return newServer(endpoints, s)
}

func (s *server) Run() {
	listener, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		log.Error(err)
		panic(err)
	}
}
