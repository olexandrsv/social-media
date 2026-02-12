package server

import (
	"social-media/internal/chats/service"
	"sync"
)

type server struct {
	httpServer Server
	grpcServer Server
}

func New(srv service.Service) Server {
	httpServer := newHttpServer(srv)
	grpcServer := newGrpcServer(srv)

	return &server{
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

func (s *server) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.httpServer.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.grpcServer.Run()
	}()
	wg.Wait()
}
