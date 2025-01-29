package transport

import (
	"context"
	"fmt"
	"net"
	"social-media/api/pb/users"
	"social-media/internal/common/app/log"
	"social-media/internal/users/endpoint"
	"sync"

	"google.golang.org/grpc"
)

type grpcServer struct {
	srv      *grpc.Server
	endpoint endpoint.Endpoints
	users.UnimplementedUsersServer
}

func newGRPCServer(e endpoint.Endpoints, s *grpc.Server) *grpcServer {
	return &grpcServer{
		endpoint: e,
		srv:      s,
	}
}

func NewGRPCServer(endpoints endpoint.Endpoints) *grpcServer {
	s := grpc.NewServer()
	server := newGRPCServer(endpoints, s)
	users.RegisterUsersServer(s, server)
	return server
}

func (s *grpcServer) UsersFullNames(ctx context.Context, req *users.UsersFullNamesReq) (*users.UsersFullNamesResp, error) {
	log.Error(fmt.Errorf("%#v", req.UserID))
	return &users.UsersFullNamesResp{
		FullName: []*users.FullName{
			&users.FullName{
				Name:    "bob",
				Surname: "smith",
			},
			&users.FullName{
				Name:    "ben",
				Surname: "arnum",
			},
		},
	}, nil
}

func (s *grpcServer) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	
	listener, err := net.Listen("tcp", ":5053")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		log.Error(err)
		panic(err)
	}
}
