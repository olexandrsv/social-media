package transport

import (
	"context"
	"net"
	"social-media/api/pb/users"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/users/domain/post"
	"social-media/internal/users/endpoint"
	"social-media/internal/users/service"
	"sync"

	"google.golang.org/grpc"
)

type grpcServer struct {
	srv      *grpc.Server
	endpoint endpoint.Endpoints
	service  service.Service
	users.UnimplementedUsersServer
}

func newGRPCServer(e endpoint.Endpoints, service service.Service, s *grpc.Server) *grpcServer {
	return &grpcServer{
		endpoint: e,
		service:  service,
		srv:      s,
	}
}

func NewGRPCServer(endpoints endpoint.Endpoints, service service.Service) *grpcServer {
	s := grpc.NewServer()
	server := newGRPCServer(endpoints, service, s)
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

func (s *grpcServer) LastReadPosts(ctx context.Context, req *users.LastReadPostsReq) (*users.LastReadPostsResp, error) {
	posts, err := s.service.LastReadPosts(int(req.UserID))
	if err != nil {
		return &users.LastReadPostsResp{
			Err: convertToError(err),
		}, nil
	}

	convertedPosts := slice.MustConvert(posts, func(post *post.Post) *users.Post {
		return &users.Post{
			Id:     post.ID,
			UserID: int64(post.UserID),
		}
	})

	return &users.LastReadPostsResp{
		Posts: convertedPosts,
	}, nil
}

func (s *grpcServer) SendPostNotification(ctx context.Context, req *users.SendPostNotificationReq) (*users.SendPostNotificationResp, error) {
	err := s.service.SendPostNotification(int(req.SenderID), req.Message)
	if err != nil {
		return &users.SendPostNotificationResp{
			Err: convertToError(err),
		}, nil
	}

	return &users.SendPostNotificationResp{
		Err: nil,
	}, nil
}

func (s *grpcServer) SendNotification(ctx context.Context, req *users.SendNotificationReq) (*users.SendNotificationResp, error) {
	receiversIDs := slice.MustConvert(req.ReceiversIDs, func(n int64) int {
		return int(n)
	})
	err := s.service.SendNotification(int(req.SenderID), receiversIDs, req.Message)
	if err != nil {
		return &users.SendNotificationResp{
			Err: convertToError(err),
		}, nil
	}

	return &users.SendNotificationResp{
		Err: nil,
	}, nil
}

func convertToError(err error) *users.Error {
	if err == nil {
		return nil
	}

	customError, ok := err.(common.Error)
	if !ok {
		customError = common.ErrInternal
	}

	return &users.Error{
		Code:    int64(customError.Code()),
		Message: customError.Message(),
	}
}

func (s *grpcServer) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", ":"+config.App.Users.Service.GrpcPort)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		log.Error(err)
		panic(err)
	}
}
