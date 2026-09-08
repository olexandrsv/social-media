package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"time"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims interface {
	jwt.Claims
	init(string, time.Time)
}

type JwtClaims struct {
	jwt.RegisteredClaims
}

func (claim *JwtClaims) init(issuer string, expiresAt time.Time) {
	claim.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
}

var jwtKey = []byte("supersercretkey")

type AuthClaim struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	*JwtClaims
}

func newAuthClaim(id int, login string) *AuthClaim {
	return &AuthClaim{
		ID:        id,
		Login:     login,
		JwtClaims: &JwtClaims{},
	}
}

type FileClaim struct {
	FileID string `json:"file_id"`
	*JwtClaims
}

func newFileClaim(fileId string) *FileClaim {
	return &FileClaim{
		FileID:    fileId,
		JwtClaims: &JwtClaims{},
	}
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
	return generateToken(1*time.Hour, newAuthClaim(id, login))
}

func (s *authService) ValidateToken(signedToken string) (int, string, error) {
	claims, err := validate(signedToken, &AuthClaim{})
	if err != nil {
		return 0, "", err
	}
	return claims.ID, claims.Login, nil
}

func (s *authService) GenerateSignedUrl(fileID string) (string, error) {
	return generateToken(5*time.Minute, newFileClaim(fileID))
}

func (s *authService) ValidateSignedUrl(signedToken string) (string, error) {
	claims, err := validate(signedToken, &FileClaim{})
	if err != nil {
		return "", err
	}
	return claims.FileID, nil
}

func generateToken(duration time.Duration, claim Claims) (string, error) {
	expiresAt := time.Now().Add(duration)
	claim.init("social-media", expiresAt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	return tokenString, nil
}

func validate[T jwt.Claims](signedToken string, claims T) (T, error) {
	var t T
	token, err := jwt.ParseWithClaims(
		signedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
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
	date, err := claims.GetExpirationTime()
	if err != nil {
		log.Error(errors.WithStack(err))
		return t, common.ErrInternal
	}

	if !date.After(time.Now()) {
		log.Error(errors.New("token expired"))
		return t, common.ErrInvalidToken
	}
	return claims, nil
}
