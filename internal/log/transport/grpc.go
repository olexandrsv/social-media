package transport

import (
	"net"
	"social-media/api/pb/log"
	"social-media/internal/log/endpoint"

	"google.golang.org/grpc"
)

type server struct {
	srv *grpc.Server
	e   endpoint.Endpoints
}

func NewGRPCServer(e endpoint.Endpoints) *server {
	s := grpc.NewServer()
	log.RegisterLogServer(s, e)
	return &server{srv: s, e: e}
}

func (s *server) Run() {
	listener, err := net.Listen("tcp", ":5052")
	if err != nil {
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		panic(err)
	}
}
