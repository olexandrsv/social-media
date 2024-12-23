package common

import (
	"context"
	"social-media/api/pb/auth"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	"github.com/pkg/errors"
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
		log.Error(errors.WithStack(err))
		panic(err)
	}
	return &client{
		auth.NewAuthenticateClient(conn),
	}
}

func (c client) GenerateToken(id int, login string) (string, error) {
	resp, err := c.GenerateJWT(context.Background(), &auth.GenerateJWTReq{Id: int64(id), Login: login})
	if err, ok := err.(Error); ok {
		return "", err
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", ErrInternal
	}
	return resp.Token, nil
}

func (c client) ValidateToken(token string) (int, string, error) {
	resp, err := c.ValidateJWT(context.Background(), &auth.ValidateJWTReq{Token: token})
	if err, ok := err.(Error); ok {
		return 0, "", err
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return 0, "", ErrInternal
	}
	return int(resp.Id), resp.Login, nil
}
