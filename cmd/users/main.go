package main

import (
	"social-media/internal/common/app"
	"social-media/internal/common/clients"
	"social-media/internal/users/endpoint"
	"social-media/internal/users/repository"
	"social-media/internal/users/service"
	"social-media/internal/users/transport"
	"sync"
)

func main() {
	app.InitUsersService()

	var wg sync.WaitGroup
	repo, err := repository.New()
	if err != nil {
		return
	}
	auth := clients.NewAuthClient()

	s := service.New(repo, auth)
	endpoints := endpoint.NewEndpoints(s)

	httpServer := transport.NewHTTPServer(s, endpoints)
	grpcServer := transport.NewGRPCServer(endpoints, s)

	wg.Add(2)
	go httpServer.Run(&wg)
	go grpcServer.Run(&wg)
	wg.Wait()
}
