package transport

import (
	"context"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/endpoint"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) postComments(w http.ResponseWriter, r *http.Request) {
	req, err := decodeCommentsReq(r)
	if err != nil {
		writeError(w, err)
		return
	}

	tone, comments, err := s.service.PostComments(req.Token, req.ParentID)
	if err != nil {
		writeError(w, err)
		return
	}

	models := slice.MustConvert(comments, commentToModel)
	writeJSON(w, PostCommentsResp{
		Tone: ToneModel{
			PositivePercentage: tone.Positive(),
			NegativePercentage: tone.Negative(),
		},
		Comments: models,
	})
}

func decodeCommentsReq(r *http.Request) (CommentsReq, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return CommentsReq{}, common.ErrNoToken
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		return CommentsReq{}, common.ErrInvalidData
	}

	return CommentsReq{
		Token:    token.Value,
		ParentID: id,
	}, nil
}

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

func (s *server) decodeDeletePostCommentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}

	params := mux.Vars(r)
	parentID, ok := params["post_id"]
	if !ok {
		return nil, common.ErrInvalidData
	}
	commentID, ok := params["comment_id"]
	if !ok {
		return nil, common.ErrInvalidData
	}

	return endpoint.DeletePostCommentReq{
		Token:     token.Value,
		ParentID:  parentID,
		CommentID: commentID,
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
		Token: token.Value,
		UpdateMessageReq: endpoint.UpdateMessageReq{
			ID:            id,
			Text:          r.FormValue("text"),
			Images:        r.MultipartForm.File["images[]"],
			Files:         r.MultipartForm.File["files[]"],
			DeletedImages: r.MultipartForm.Value["deletedImages[]"],
			DeletedFiles:  r.MultipartForm.Value["deletedFiles[]"],
		},
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

func (s *server) decodeCreateCommnetCommentReq(ctx context.Context, r *http.Request) (interface{}, error) {
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

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.CreateCommentCommentReq{
		Token:    token.Value,
		ParentID: id,
		Text:     r.FormValue("text"),
		Images:   r.MultipartForm.File["images[]"],
		Files:    r.MultipartForm.File["files[]"],
	}, nil
}

func (s *server) decodeDeleteCommentCommentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}

	params := mux.Vars(r)
	parentID, ok := params["parent_id"]
	if !ok {
		return nil, common.ErrInvalidData
	}
	commentID, ok := params["comment_id"]
	if !ok {
		return nil, common.ErrInvalidData
	}

	return endpoint.DeleteCommentCommentReq{
		Token:     token.Value,
		ParentID:  parentID,
		CommentID: commentID,
	}, nil
}
