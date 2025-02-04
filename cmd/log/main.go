package main

import (
	"social-media/internal/common/app"
	"social-media/internal/log/endpoint"
	"social-media/internal/log/service"
	"social-media/internal/log/transport"
)

func main() {
	app.InitLogService()
	
	s := service.NewService()
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewGRPCServer(endpoints)
	server.Run()
}
