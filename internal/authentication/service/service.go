package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersercretkey")

type JWTClaim struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}

type Service interface {
	GenerateToken(int, string) (string, error)
	ValidateToken(string) (int, string, error)
}

type authService struct {
}

func New() Service {
	return &authService{}
}

func (s *authService) GenerateToken(id int, login string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		ID:    id,
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	return tokenString, nil
}

func (s *authService) ValidateToken(signedToken string) (int, string, error) {
	claims, err := parseToken(signedToken)
	if err != nil {
		log.Error(err)
		return 0, "", common.ErrInvalidToken
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Error(errors.New("token expired"))
		return claims.ID, claims.Login, common.ErrInvalidToken
	}
	return claims.ID, claims.Login, nil
}

func parseToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return claims, errors.New("couldn't parse claims")
	}
	return claims, nil
}
