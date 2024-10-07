package main

import (
	"social-media/internal/common"
	"social-media/internal/common/app"
	"social-media/internal/users/endpoint"
	"social-media/internal/users/repository"
	"social-media/internal/users/service"
	"social-media/internal/users/transport"
)

func main() {
	app.InitUsersService()
	repo := repository.New()
	auth := common.NewAuthClient()

	s := service.New(repo, auth)
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewHTTPServer(endpoints)

	server.Run()
}
