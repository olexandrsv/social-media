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
	OwnPosts(ctx context.Context, request interface{}) (interface{}, error)
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
	return Post{
		ID:          post.ID(),
		Text:        post.Text(),
		FilesPaths:  post.FilesPaths(),
		ImagesPaths: post.ImagesPaths(),
	}, nil
}

func (e *postsEndpoint) OwnPosts(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(Token)
	if !ok {
		log.Error(errors.New("can't assign to Token"))
		return nil, common.ErrInternal
	}

	posts, err := e.s.OwnPosts(req.Token)
	if err != nil {
		return nil, err
	}

	postModels := make([]Post, 0, len(posts))
	for _, post := range posts {
		postModels = append(postModels, Post{
			ID:          post.ID(),
			Text:        post.Text(),
			FilesPaths:  post.FilesPaths(),
			ImagesPaths: post.ImagesPaths(),
		})
	}
	return postModels, nil
}