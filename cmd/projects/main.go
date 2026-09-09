package main

import (
	"social-media/internal/common/app"
	"social-media/internal/common/clients"
	"social-media/internal/projects/repository"
	"social-media/internal/projects/server"
	"social-media/internal/projects/service"
)

func main() {
	app.InitProjectsService()

	auth := clients.NewAuthClient()
	users := clients.NewUsersClient()

	repo, err := repository.New()
	if err != nil {
		return
	}
	s := service.New(repo, auth, users)
	srv := server.New(s)
	srv.Run()
}
