package main

import (
	"social-media/internal/chats/repository"
	"social-media/internal/chats/server"
	"social-media/internal/chats/service"
	"social-media/internal/common/app"
	"social-media/internal/common/clients"
)

func main() {
	app.InitChatsService()
	
	repo, err := repository.New()
	if err != nil {
		return
	}

	authClient := clients.NewAuthClient()
	usersClient := clients.NewUsersClient()

	service := service.New(repo, authClient, usersClient)
	srv := server.New(service)
	srv.Run()
}