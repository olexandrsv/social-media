package clients

import (
	"context"
	"social-media/api/pb/auth"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type AuthClient interface {
	GenerateToken(int, string) (string, error)
	ValidateToken(string) (int, string, error)
	GenerateSignedUrl(string) (string, error)
	ValidateSignedUrl(string) (string, error)
}

type client struct {
	c auth.AuthenticateClient
}

func NewAuthClient() AuthClient {
	host := config.App.AuthService.Host
	port := config.App.AuthService.Port
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(err)
	}
	return &client{
		auth.NewAuthenticateClient(conn),
	}
}

func (c client) GenerateToken(id int, login string) (string, error) {
	resp, err := c.c.GenerateJWT(context.Background(), &auth.GenerateJWTReq{Id: int64(id), Login: login})
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	if resp.Err != nil {
		return "", common.NewError(int(resp.Err.Code), resp.Err.Message)
	}
	return resp.Token, nil
}

func (c client) ValidateToken(token string) (int, string, error) {
	resp, err := c.c.ValidateJWT(context.Background(), &auth.ValidateJWTReq{Token: token})
	if err != nil {
		log.Error(errors.WithStack(err))
		return 0, "", common.ErrInternal
	}
	if resp.Err != nil {
		return 0, "", common.NewError(int(resp.Err.Code), resp.Err.Message)
	}
	return int(resp.Id), resp.Login, nil
}

func (c client) GenerateSignedUrl(fileID string) (string, error) {
	resp, err := c.c.GenerateSignedUrl(context.Background(), &auth.GenerateSignedUrlReq{
		FileID: fileID,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	if resp.Err != nil {
		return "", common.NewError(int(resp.Err.Code), resp.Err.Message)
	}

	return resp.Token, nil
}

func (c client) ValidateSignedUrl(token string) (string, error) {
	resp, err := c.c.ValidateSignedUrl(context.Background(), &auth.ValidateSignedUrlReq{Token: token})
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	if resp.Err != nil {
		return "", common.NewError(int(resp.Err.Code), resp.Err.Message)
	}

	return resp.FileID, nil
}
