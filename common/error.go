package common

import (
	"net/http"
)

var (
	ErrInternal         = NewError(http.StatusInternalServerError, "Internal error")
	ErrInvalidData      = NewError(http.StatusBadRequest, "Invalid data")
	ErrWrongCredentials = NewError(http.StatusUnauthorized, "Invalid login or password")
	ErrNoLogin          = NewError(http.StatusBadRequest, "Login not present")
	ErrNoToken          = NewError(http.StatusUnauthorized, "Token not present")
	ErrNotFound         = NewError(http.StatusNotFound, "No Data")
	ErrLoginExists      = NewError(http.StatusOK, "Login exists")

	ErrInvalidToken = NewError(http.StatusUnauthorized, "Invalid Token")
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}
