package endpoint

import (
	"context"
	"errors"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/service"
)

type Endpoint interface {
	CreatePost(ctx context.Context, request interface{}) (interface{}, error)
	GetPosts(ctx context.Context, request interface{}) (interface{}, error)
	UpdatePost(ctx context.Context, request interface{}) (interface{}, error)
	DeletePost(ctx context.Context, request interface{}) (interface{}, error)
	PostComments(ctx context.Context, request interface{}) (interface{}, error)
	commentEndpoint
}

type postsEndpoint struct {
	s service.Service
}

func New(s service.Service) Endpoint {
	return &postsEndpoint{
		s: s,
	}
}

func (e *postsEndpoint) CreatePost(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(CreatePostReq)
	if !ok {
		log.Error(errors.New("can't assign to CreatePostReq"))
		return nil, common.ErrInternal
	}

	post, err := e.s.CreatePost(req.Token, req.Text, req.FilesPath, req.ImagesPath)
	if err != nil {
		return nil, err
	}
	return PostModel{
		ID:          post.ID(),
		Text:        post.Text(),
		FilesPaths:  post.FilesPaths(),
		ImagesPaths: post.ImagesPaths(),
	}, nil
}

func (e *postsEndpoint) GetPosts(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(GetPostsRequest)
	if !ok {
		log.Error(errors.New("can't assign to GetPostsRequest"))
		return nil, common.ErrInternal
	}

	posts, err := e.s.GetPosts(req.Token, req.UserID)
	if err != nil {
		return nil, err
	}

	postModels := make([]PostModel, 0, len(posts))
	for _, post := range posts {
		postModels = append(postModels, PostModel{
			ID:          post.ID(),
			UserID:      post.UserID(),
			Text:        post.Text(),
			FilesPaths:  post.FilesPaths(),
			ImagesPaths: post.ImagesPaths(),
		})
	}
	return postModels, nil
}

func (e *postsEndpoint) UpdatePost(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(UpdatePostReq)
	if !ok {
		log.Error(errors.New("can't assign to UpdatePostReq"))
		return nil, common.ErrInternal
	}

	post, err := e.s.UpdatePost(service.UpdatePostReq{
		Token:         req.Token,
		ID:            req.ID,
		Text:          req.Text,
		Images:        req.Images,
		Files:         req.Files,
		DeletedImages: req.DeletedImages,
		DeletedFiles:  req.DeletedFiles,
	})
	if err != nil {
		return nil, err
	}

	return PostModel{
		ID:          post.ID(),
		UserID:      post.UserID(),
		Text:        post.Text(),
		FilesPaths:  post.FilesPaths(),
		ImagesPaths: post.ImagesPaths(),
	}, nil
}

func (e *postsEndpoint) DeletePost(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(DeletePostReq)
	if !ok {
		log.Error(errors.New("can't assign to DeletePostReq"))
		return nil, common.ErrInternal
	}

	err := e.s.DeletePost(req.Token, req.PostID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *postsEndpoint) PostComments(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(GetPostCommentsReq)
	if !ok {
		log.Error(errors.New("can't assign to GetPostCommentsReq"))
		return nil, common.ErrInternal
	}

	comments, err := e.s.PostComments(req.Token, req.PostID)
	if err != nil {
		return nil, err
	}

	commentsModels := make([]CommentModel, 0, len(comments))
	for _, comment := range comments {
		commentModel := CommentModel{
			ID:         comment.ID(),
			UserID:     comment.UserID(),
			Text:       comment.Text(),
			ImagesPath: comment.ImagesPaths(),
			FilesPath:  comment.FilesPaths(),
		}
		commentsModels = append(commentsModels, commentModel)
	}
	return commentsModels, nil
}
