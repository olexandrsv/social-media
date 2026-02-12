package main

import (
	"social-media/internal/common/app"
	"social-media/internal/log/server"
	"social-media/internal/log/service"
)

func main() {
	app.InitLogService()

	s := service.NewService()
	server := server.NewServer(s)
	server.Run()
}
