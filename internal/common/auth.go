package common

import (
	"context"
	"log"
	"social-media/api/pb/auth"
	"social-media/internal/common/app/config"

	"google.golang.org/grpc"
)

type AuthClient interface {
	GenerateToken(int, string) (string, error)
	ValidateToken(string) (int, string, error)
}

type client struct {
	auth.AuthenticateClient
}

func NewAuthClient() AuthClient {
	conn, err := grpc.Dial(":"+config.App.AuthService.Port, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil
	}
	return &client{
		auth.NewAuthenticateClient(conn),
	}
}

func (c client) GenerateToken(id int, login string) (string, error) {
	resp, err := c.GenerateJWT(context.Background(), &auth.GenerateJWTReq{Id: int64(id), Login: login})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c client) ValidateToken(token string) (int, string, error) {
	resp, err := c.ValidateJWT(context.Background(), &auth.ValidateJWTReq{Token: token})
	if err != nil {
		return 0, "", err
	}
	return int(resp.Id), resp.Login, nil
}
