package main

import (
	"social-media/internal/common"
	"social-media/internal/common/app"
	"social-media/internal/posts/endpoint"
	"social-media/internal/posts/repository"
	"social-media/internal/posts/service"
	"social-media/internal/posts/transport"
)

func main() {
	app.InitPostsService()

	repo := repository.New()
	auth := common.NewAuthClient()

	s := service.New(repo, auth)
	e := endpoint.New(s)

	r := transport.NewHTTPServer(e)
	r.Run()
}
