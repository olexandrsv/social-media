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
	ErrForbidden          = NewError(http.StatusForbidden, "Forbidden")

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

// type CustomError interface {
// 	OriginalErr() error
// 	Error() string
// }

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

// type tracesSkiper struct {
// 	err  error
// 	skip int
// }

// func newTracesSkiper(err error, skip int) tracesSkiper {
// 	return tracesSkiper{
// 		err:  err,
// 		skip: skip,
// 	}
// }

// func (ts tracesSkiper) Error() string {
// 	tracer, ok := ts.err.(stackTracer)
// 	if !ok {
// 		return fmt.Sprintf("%+v\n", ts.err)
// 	}
// 	trace := tracer.StackTrace()
// 	if len(trace) < ts.skip {
// 		return fmt.Sprintf("%+v\n", ts.err)
// 	}
// 	trace = trace[ts.skip:]
// 	return fmt.Sprintf("%+v\n", trace)
// }

// type customError struct {
// 	errorMsg    string
// 	originalErr tracesSkiper
// }

// func NewCustomError(err error, errorMsg string) CustomError {
// 	return customError{
// 		errorMsg: errorMsg,
// 		originalErr: newTracesSkiper(errors.New(err.Error()), 2),
// 	}
// }

// func (err customError) Error() string {
// 	return err.errorMsg
// }

// func (err customError) OriginalErr() error {
// 	return err.originalErr
// }

// type NotFoundError struct {
// 	CustomError
// }

// func NewNotFoundError(err error, title string, id int) NotFoundError {
// 	return NotFoundError{
// 		CustomError: NewCustomError(err, fmt.Sprintf("item '%s' with id %d wasn't found", title, id)),
// 	}
// }

// type InternalError struct{
// 	CustomError
// }

// func NewIternalError(err error) InternalError{
// 	return InternalError{
// 		CustomError: NewCustomError(err, "Internal error"),
// 	}
// }

// type InvalidDataError struct{
// 	CustomError
// }

// func NewInvalidDataError(err error) InvalidDataError{
// 	return InvalidDataError{
// 		CustomError: NewCustomError(err, "Invalid data"),
// 	}
// }
