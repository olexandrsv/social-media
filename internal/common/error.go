package common

import (
	"net/http"
)

var (
	ErrInternal           = NewError(http.StatusInternalServerError, "Internal error")
	ErrInvalidData        = NewError(http.StatusBadRequest, "Invalid data")
	ErrWrongCredentials   = NewError(http.StatusUnauthorized, "Invalid login or password")
	ErrNoToken            = NewError(http.StatusUnauthorized, "Token not present")
	ErrNotFound           = NewError(http.StatusNotFound, "No Data")
	ErrLoginExists        = NewError(http.StatusOK, "Login exists")
	ErrSubscriptionExists = NewError(http.StatusBadRequest, "Subscription exists")

	ErrInvalidToken = NewError(http.StatusUnauthorized, "Invalid Token")
)

type Error struct {
	code    int
	message string
}

func NewError(code int, msg string) Error {
	return Error{
		code:    code,
		message: msg,
	}
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Error() string {
	return e.message
}
