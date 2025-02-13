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
	DeletePostComment(ctx context.Context, request interface{}) (interface{}, error)
	UpdateComment(ctx context.Context, request interface{}) (interface{}, error)
	CommentComments(ctx context.Context, request interface{}) (interface{}, error)
	CreateCommentComment(ctx context.Context, request interface{}) (interface{}, error)
	DeleteCommentComment(ctx context.Context, request interface{}) (interface{}, error)
}

func (e *postsEndpoint) CreatePostComment(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CreatePostCommentReq)
	if !ok {
		log.Error(errors.New("can't assign request to CreatePostCommentReq"))
		return nil, common.ErrInternal
	}

	comment, err := e.s.CreatePostComment(service.CreatePostCommentReq{
		PostID: req.PostID,
		CreateMessageReq: service.CreateMessageReq{
			Token:  req.Token,
			Text:   req.Text,
			Images: req.Images,
			Files:  req.Files,
		},
	})
	if err != nil {
		return nil, err
	}

	return commentToModel(comment), nil
}

func (e *postsEndpoint) UpdateComment(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(UpdateCommentReq)
	if !ok {
		log.Error(errors.New("can't assign request to UpdateCommentReq"))
		return nil, common.ErrInternal
	}

	comment, err := e.s.UpdateComment(service.UpdateCommentReq{
		UpdateMessageReq: service.UpdateMessageReq{
			Token:         req.Token,
			ID:            req.ID,
			Text:          req.Text,
			Images:        req.Images,
			Files:         req.Files,
			DeletedImages: req.DeletedImages,
			DeletedFiles:  req.DeletedFiles,
		},
	})
	if err != nil {
		return nil, err
	}

	return commentToModel(comment), nil
}

func (e *postsEndpoint) CommentComments(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(GetCommentCommentsReq)
	if !ok {
		log.Error(errors.New("can't assign to GetCommentCommentsReq"))
		return nil, common.ErrInternal
	}

	comments, err := e.s.CommentComments(req.Token, req.CommentID)
	if err != nil {
		return nil, err
	}

	return commentsToModels(comments), nil
}

func (e *postsEndpoint) CreateCommentComment(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CreateCommentCommentReq)
	if !ok {
		log.Error(errors.New("can't assign to CreateCommentCommentReq"))
		return nil, common.ErrInternal
	}

	comment, err := e.s.CreateCommentComment(service.CreateCommentCommentReq{
		ParentID: req.ParentID,
		CreateMessageReq: service.CreateMessageReq{
			Token:  req.Token,
			Text:   req.Text,
			Images: req.Images,
			Files:  req.Files,
		},
	})
	if err != nil {
		return nil, err
	}

	return commentToModel(comment), nil
}

func (e *postsEndpoint) DeletePostComment(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(DeletePostCommentReq)
	if !ok {
		log.Error(errors.New("can't assign to DeletePostCommentReq"))
		return nil, common.ErrInternal
	}

	err := e.s.DeletePostComment(req.Token, req.ParentID, req.CommentID)
	return nil, err
}

func (e *postsEndpoint) DeleteCommentComment(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(DeleteCommentCommentReq)
	if !ok {
		log.Error(errors.New("can't assign to DeleteCommentCommentReq"))
		return nil, common.ErrInternal
	}

	err := e.s.DeleteCommentComment(req.Token, req.ParentID, req.CommentID)
	return nil, err
}
