package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

type Claims interface {
	Valid() error
	ExpiresAt() int64
	SetExpirationTime(int64)
}

type StandardClaim struct {
	jwt.StandardClaims
}

func (claim *StandardClaim) ExpiresAt() int64 {
	return claim.StandardClaims.ExpiresAt
}

func (claim *StandardClaim) SetExpirationTime(time int64) {
	claim.StandardClaims.ExpiresAt = time
}

var jwtKey = []byte("supersercretkey")

type JWTClaim struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	*StandardClaim
}

type FileClaim struct {
	FileID string `json:"file_id"`
	*StandardClaim
}

func (claim *FileClaim) ExpiresAt() int64 {
	return claim.StandardClaims.ExpiresAt
}

type Service interface {
	GenerateToken(int, string) (string, error)
	ValidateToken(string) (int, string, error)
	GenerateSignedUrl(string) (string, error)
	ValidateSignedUrl(string) (string, error)
}

type authService struct {
}

func New() Service {
	return &authService{}
}

func (s *authService) GenerateToken(id int, login string) (string, error) {
	return generateToken(1*time.Hour, &JWTClaim{
		ID:            id,
		Login:         login,
		StandardClaim: &StandardClaim{},
	})
}

func (s *authService) ValidateToken(signedToken string) (int, string, error) {
	claims, err := validate(signedToken, &JWTClaim{})
	if err != nil {
		return 0, "", err
	}
	return claims.ID, claims.Login, nil
}

func (s *authService) GenerateSignedUrl(fileID string) (string, error) {
	return generateToken(5*time.Minute, &FileClaim{
		FileID:        fileID,
		StandardClaim: &StandardClaim{},
	})
}

func (s *authService) ValidateSignedUrl(signedToken string) (string, error) {
	claims, err := validate(signedToken, &FileClaim{})
	if err != nil {
		return "", err
	}
	return claims.FileID, nil
}

func generateToken(duration time.Duration, claim Claims) (string, error) {
	expirationTime := time.Now().Add(duration)
	claim.SetExpirationTime(expirationTime.Unix())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	return tokenString, nil
}

func validate[T Claims](signedToken string, claims T) (T, error) {
	var t T
	token, err := jwt.ParseWithClaims(
		signedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		log.Error(errors.WithStack(err))
		return t, common.ErrInvalidToken
	}
	claims, ok := token.Claims.(T)
	if !ok {
		log.Error(errors.New("couldn't parse claims"))
		return t, common.ErrInvalidToken
	}

	if claims.ExpiresAt() < time.Now().Local().Unix() {
		log.Error(errors.New("token expired"))
		return t, common.ErrInvalidToken
	}
	return claims, nil
}
