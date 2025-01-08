package endpoint

import (
	"context"
	"errors"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/service"
)

type commentEndpoint interface {
	CreatePostComment(ctx context.Context, request interface{}) (interface{}, error)
}

func (e *postsEndpoint) CreatePostComment(ctx context.Context, request interface{}) (interface{}, error){
	req, ok := request.(CreatePostCommentReq)
	if !ok {
		log.Error(errors.New("can't assign request to CreatePostCommentReq"))
		return nil, common.ErrInternal
	}

	comment, err := e.s.CreatePostComment(service.CreatePostCommentReq{
		Token: req.Token,
		PostID: req.PostID,
		Text: req.Text,
		Images: req.Images,
		Files: req.Files,
	})
	if err != nil {
		return nil, err
	}

	return CommentModel{
		ID: comment.ID(),
		UserID: comment.UserID(),
		Text: comment.Text(),
		ImagesPath: comment.ImagesPaths(),
		FilesPath: comment.FilesPaths(),
	}, nil
}