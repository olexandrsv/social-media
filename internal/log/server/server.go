package server

import (
	"social-media/internal/log/service"
	"sync"
)

type server struct {
	grpcServer *grpcServer
	httpServer *httpServer
}

func NewServer(service service.Service) *server {
	return &server{
		grpcServer: newGrpcServer(service),
		httpServer: newHttpServer(service),
	}
}

func (s *server) Run() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		s.httpServer.run()
	}()
	go func() {
		defer wg.Done()
		s.grpcServer.run()
	}()

	wg.Wait()
}
