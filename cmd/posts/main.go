package main

import (
	"social-media/internal/common/app"
	"social-media/internal/common/clients"
	"social-media/internal/posts/endpoint"
	"social-media/internal/posts/repository"
	"social-media/internal/posts/service"
	"social-media/internal/posts/transport"
)

func main() {
	app.InitPostsService()

	repo := repository.New()
	auth := clients.NewAuthClient()
	users := clients.NewUsersClient()

	s := service.New(repo, auth, users)
	e := endpoint.New(s)

	r := transport.NewHTTPServer(e)
	r.Run()
}
