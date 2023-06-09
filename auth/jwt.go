package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersercretkey")

type JWTClaim struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	jwt.StandardClaims
}

func GenerateJWT(id int, login string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		ID:    id,
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	claims, err := parseToken(signedToken)
	if err != nil {
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func TokenCredentials(token string) (id int, login string, err error) {
	claims, err := parseToken(token)
	if err != nil {
		return
	}
	return claims.ID, claims.Login, err
}

func parseToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	return claims, nil
}
