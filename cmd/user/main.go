package main

import (
	"log"
	"os"
	"social-media/common"
	"social-media/users"
	"social-media/users/endpoint"
	"social-media/users/transport"
)

const url = `postgresql://postgres:database@localhost:5432/media?sslmode=disable`

func main() {
	repo := users.NewRepository(url)
	logger := NewLogger("../../log.txt")
	auth := common.NewAuthClient()

	s := users.NewService(repo, logger, auth)
	endpoints := endpoint.NewEndpoints(s)
	server := transport.NewHTTPServer(endpoints)

	server.Run()
}

func NewLogger(path string) *log.Logger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal("Failed to open log file")
	}

	return log.New(file, "User microservice: ", log.Ldate|log.Ltime|log.Lshortfile)
}
