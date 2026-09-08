package main

import (
	"social-media/internal/authentication/service"
	"social-media/internal/authentication/server"
	"social-media/internal/common/app"
)

func main() {
	app.InitAuthService()

	s, err := service.New()
	if err != nil {
		return
	}
	server := server.NewGRPCServer(s)
	server.Run()
}
