package main

import (
	"social-media/internal/authentication/endpoint"
	"social-media/internal/authentication/service"
	"social-media/internal/authentication/transport"
)

func main() {
	s := service.New()
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewGRPCServer(endpoints)
	server.Run()
}
