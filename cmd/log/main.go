package main

import (
	"social-media/internal/common/app"
	"social-media/internal/log/service"
	"social-media/internal/log/server"
)

func main() {
	app.InitLogService()

	s := service.NewService()
	server := server.NewGRPCServer(s)
	server.Run()
}
