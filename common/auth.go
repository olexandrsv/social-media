package common

import (
	"context"
	"log"
	"social-media/api/pb"
	"google.golang.org/grpc"
)

type AuthClient interface{
	GenerateToken(int, string) (string, error)
	ValidateToken(string) (int, string, error)
}

type client struct{
	pb.AuthenticateClient
}

func NewAuthClient() AuthClient {
	conn, err := grpc.Dial(":5051", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil
	}
	return &client{
		pb.NewAuthenticateClient(conn),
	}
}

func (c client) GenerateToken(id int, login string) (string, error) {
	resp, err := c.GenerateJWT(context.Background(), &pb.GenerateJWTReq{Id: int64(id), Login: login})
	if err != nil{
		return "", err
	}
	return resp.Token, nil
}

func (c client) ValidateToken(token string) (int, string, error) {
	resp, err := c.ValidateJWT(context.Background(), &pb.ValidateJWTReq{Token: token})
	if err != nil {
		return 0, "", err
	}
	return int(resp.Id), resp.Login, nil
}