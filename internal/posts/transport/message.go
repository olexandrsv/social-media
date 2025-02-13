package transport

import (
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/service"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) decodeCreateMessageReq(r *http.Request) (service.CreateMessageReq, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return service.CreateMessageReq{}, common.ErrNoToken
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return service.CreateMessageReq{}, common.ErrInvalidData
	}

	return service.CreateMessageReq{
		Token:  token.Value,
		Text:   r.FormValue("text"),
		Images: r.MultipartForm.File["images[]"],
		Files:  r.MultipartForm.File["files[]"],
	}, nil
}

func (s *server) decodeUpdateMessageReq(r *http.Request) (service.UpdateMessageReq, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return service.UpdateMessageReq{}, common.ErrNoToken
	}

	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		return service.UpdateMessageReq{}, common.ErrInvalidData
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return service.UpdateMessageReq{}, common.ErrInvalidData
	}

	return service.UpdateMessageReq{
		Token:         token.Value,
		ID:            id,
		Text:          r.FormValue("text"),
		Images:        r.MultipartForm.File["images[]"],
		Files:         r.MultipartForm.File["files[]"],
		DeletedImages: r.MultipartForm.Value["deletedImages[]"],
		DeletedFiles:  r.MultipartForm.Value["deletedFiles[]"],
	}, nil
}
