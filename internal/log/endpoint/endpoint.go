package endpoint

import (
	"context"
	"social-media/api/pb/log"
	"social-media/internal/log/service"
)

type Endpoints struct {
	log.UnimplementedLogServer
	service service.Service
}

func NewEndpoints(service service.Service) Endpoints{
	return Endpoints{service: service}
}

func (e Endpoints) Error(ctx context.Context, req *log.LogRequest) (*log.Empty, error){
	e.service.Error(req.Msg)
	return &log.Empty{}, nil
}