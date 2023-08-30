package main

import (
	"log"
	"net"
	"social-media/api/pb"
	"social-media/authentication"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":5051")
	if err != nil{
		log.Println(err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthenticateServer(s, authentication.NewServer())

	if err := s.Serve(listener); err != nil{
		log.Println(err)
	}
}