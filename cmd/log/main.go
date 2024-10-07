package main

import (
	"social-media/internal/log/endpoint"
	"social-media/internal/log/service"
	"social-media/internal/log/transport"
)

func main() {
	s := service.NewService()
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewGRPCServer(endpoints)
	server.Run()
}
