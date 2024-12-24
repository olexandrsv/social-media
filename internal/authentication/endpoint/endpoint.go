package endpoint

import (
	"context"
	"social-media/api/pb/auth"
	"social-media/internal/authentication/service"
	"social-media/internal/common"
)

type Endpoints struct {
	auth.UnimplementedAuthenticateServer
	service service.Service
}

func NewEndpoints(s service.Service) Endpoints {
	return Endpoints{service: s}
}

func (e Endpoints) GenerateJWT(ctx context.Context, req *auth.GenerateJWTReq) (*auth.GenerateJWTResp, error) {
	token, err := e.service.GenerateToken(int(req.Id), req.Login)

	return &auth.GenerateJWTResp{
		Token: token,
		Err:   convertToError(err),
	}, nil
}

func (e Endpoints) ValidateJWT(ctx context.Context, req *auth.ValidateJWTReq) (*auth.ValidateJWTResp, error) {
	id, login, err := e.service.ValidateToken(req.Token)

	return &auth.ValidateJWTResp{
		Id:    int64(id),
		Login: login,
		Err:   convertToError(err),
	}, nil
}

func convertToError(err error) *auth.Error {
	if err == nil {
		return nil
	}

	customError, ok := err.(common.Error)
	if !ok {
		customError = common.ErrInternal
	}

	return &auth.Error{
		Code:    int64(customError.Code()),
		Message: customError.Message(),
	}
}
