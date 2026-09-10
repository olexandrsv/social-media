package main

import (
	"social-media/internal/common/app"
	"social-media/internal/common/clients"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/endpoint"
	"social-media/internal/posts/repository"
	"social-media/internal/posts/service"
	"social-media/internal/posts/transport"
)

func main() {
	app.InitPostsService()

	repo, err := repository.New()
	if err != nil {
		return
	}
	auth := clients.NewAuthClient()
	users := clients.NewUsersClient()
	chats := clients.NewChatsClient()
	ai := clients.NewAIClient[*comment.Comment]()

	s := service.New(repo, auth, users, chats, ai)
	e := endpoint.New(s)

	r := transport.NewHTTPServer(e, s)
	r.Run()
}
