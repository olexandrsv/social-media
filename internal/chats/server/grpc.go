package server

import (
	"context"
	"net"
	"social-media/api/pb/chats"
	"social-media/internal/chats/service"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"

	"google.golang.org/grpc"
)

type grpcServer struct {
	srv     *grpc.Server
	service service.Service
	chats.UnimplementedChatsServer
}

func newGrpcServer(service service.Service) Server {
	s := grpc.NewServer()
	srv := &grpcServer{
		srv:     s,
		service: service,
	}
	chats.RegisterChatsServer(s, srv)
	return srv
}

func (s *grpcServer) LastReadMessages(ctx context.Context, req *chats.LastReadMessagesReq) (*chats.LastReadMessagesResp, error) {
	messages, err := s.service.LastReadMessages(int(req.UserID))
	if err != nil {
		return &chats.LastReadMessagesResp{
			Messages: nil,
			Err:      convertToError(err),
		}, nil
	}
	messageModels := slice.MustConvert(messages, func(m *chatmessage.ChatMessage) *chats.Message {
		return &chats.Message{
			LastMessageID: m.ID(),
			ChatID:        int64(m.ChatID()),
		}
	})
	return &chats.LastReadMessagesResp{
		Messages: messageModels,
	}, nil
}

func convertToError(err error) *chats.Error {
	if err == nil {
		return nil
	}

	customError, ok := err.(common.Error)
	if !ok {
		customError = common.ErrInternal
	}

	return &chats.Error{
		Code:    int64(customError.Code()),
		Message: customError.Message(),
	}
}

func (s *grpcServer) Run() {
	listener, err := net.Listen("tcp", ":"+config.App.Chats.Service.GrpcPort)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := s.srv.Serve(listener); err != nil {
		log.Error(err)
		panic(err)
	}
}
