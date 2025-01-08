package transport

import (
	"context"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/endpoint"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) decodeCreatePostCommentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}

	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		return nil, common.ErrInvalidData
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.CreatePostCommentReq{
		Token:  token.Value,
		PostID: id,
		Text:   r.FormValue("text"),
		Images: r.MultipartForm.File["images[]"],
		Files:  r.MultipartForm.File["files[]"],
	}, nil
}

func (s *server) decodeUpdateCommentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}

	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		return nil, common.ErrInvalidData
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.UpdateCommentReq{
		Token:         token.Value,
		CommentID:     id,
		Text:          r.FormValue("text"),
		Images:        r.MultipartForm.File["images[]"],
		Files:         r.MultipartForm.File["files[]"],
		DeletedImages: r.MultipartForm.Value["deletedImages[]"],
		DeletedFiles:  r.MultipartForm.Value["deletedFiles[]"],
	}, nil
}

func (s *server) decodeCommentCommentsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.GetCommentCommentsReq{
		Token:     token.Value,
		CommentID: id,
	}, nil
}
