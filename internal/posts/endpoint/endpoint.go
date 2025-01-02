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
