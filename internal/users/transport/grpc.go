package transport

import (
	"context"
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

func (s *grpcServer) UsersInfo(ctx context.Context, req *users.UsersInfoReq) (*users.UsersInfoResp, error) {
	convertedIDs := make([]int, 0, len(req.UserID))
	for _, id := range req.UserID {
		convertedIDs = append(convertedIDs, int(id))
	}
	resp, err := s.endpoint.UsersInfo(ctx, endpoint.UsersInfoReq{
		IDs: convertedIDs,
	})
	if err != nil {
		return nil, err
	}

	infos := make([]*users.Info, 0, len(resp.Info))
	for _, info := range resp.Info {
		infos = append(infos, &users.Info{
			Login:   info.Login,
			Name:    info.Name,
			Surname: info.Surname,
		})
	}

	return &users.UsersInfoResp{
		Info: infos,
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
