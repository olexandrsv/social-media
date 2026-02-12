package common

import (
	"encoding/json"
	"io"
	"net/http"
	"social-media/internal/common/app/log"

	"github.com/pkg/errors"
)

func WriteResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, resp)
}

func WriteError(w http.ResponseWriter, err error) {
	code := 500
	msg := "Internal server error"
	if e, ok := err.(Error); ok {
		code = e.Code()
		msg = e.Message()
	}
	w.WriteHeader(code)

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}

func WriteJSON(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}

func ParseResponse[T any](resp *http.Response) (T, error) {
	var t T
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return t, ErrInternal
	}

	json.Unmarshal(content, &t)

	if resp.StatusCode == 200 {
		return t, nil
	}
	return t, NewError(resp.StatusCode, string(content))
}

func NewHealthHandler(serverName string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(serverName+" is running"))
		if err != nil {
			log.Error(err)
		}
	}
}
