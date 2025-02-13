package transport

import (
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/service"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) updatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	post, err := s.service.UpdatePost(service.UpdatePostReq{
		UpdateMessageReq: service.UpdateMessageReq{
			Token:         token.Value,
			ID:            id,
			Text:          r.FormValue("text"),
			Images:        r.MultipartForm.File["images[]"],
			Files:         r.MultipartForm.File["files[]"],
			DeletedImages: r.MultipartForm.Value["deletedImages[]"],
			DeletedFiles:  r.MultipartForm.Value["deletedFiles[]"],
		},
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, UpdatePostResp(postToModel(post)))
}

func (s *server) createPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	req := service.CreatePostReq{
		CreateMessageReq: service.CreateMessageReq{
			Token:  token.Value,
			Text:   r.FormValue("text"),
			Images: r.MultipartForm.File["images[]"],
			Files:  r.MultipartForm.File["files[]"],
		},
	}

	post, err := s.service.CreatePost(req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, CreatePostResp(postToModel(post)))
}

func (s *server) getPosts(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}
	params := mux.Vars(r)
	rawID, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	posts, err := s.service.GetPosts(token.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	models := slice.MustConvert(posts, postToModel)

	writeJSON(w, GetPostsResp(models))
}

func (s *server) deletePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	err = s.service.DeletePost(token.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, nil)
}
