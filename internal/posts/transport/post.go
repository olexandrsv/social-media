package transport

import (
	"fmt"
	"io"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/post"
	"social-media/internal/posts/models"
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

	writeJSON(w, UpdatePostResp(models.PostToModel(post)))
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

	writeJSON(w, CreatePostResp(models.PostToModel(post)))
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
		log.Error(err)
		writeError(w, err)
		return
	}

	models := slice.MustConvert(posts, models.PostToModel)

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

func (s *server) getMissedPosts(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}

	missedPostsNumber, err := s.service.MissedPostsNumber(token.Value)
	if err != nil {
		writeError(w, err)
		return
	}

	converted := slice.MustConvert(missedPostsNumber, func(n *post.MissedNumber) MissedPostsNumber {
		return MissedPostsNumber{
			FollowingID: n.FollowingID(),
			Number:      n.Number(),
		}
	})
	writeJSON(w, converted)
}

func (s *server) getPostFile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	id, ok := mux.Vars(r)["file_id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}

	file, err := s.service.GetPostFile(cookie.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name()))
	w.Header().Set("Content-Type", "text/plain")

	_, err = io.Copy(w, file)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}
