package main

import (
	"log"
	"os"
	"social-media/internal/authentication/service"
	"social-media/internal/authentication/endpoint"
	"social-media/internal/authentication/transport"
)

func main() {
	logger := NewLogger("../../log.txt")
	s := service.NewService(logger)
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewGRPCServer(endpoints)
	server.Run()
}

func NewLogger(path string) *log.Logger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal("Failed to open log file")
	}

	return log.New(file, "Auth microservice: ", log.Ldate|log.Ltime|log.Lshortfile)
}