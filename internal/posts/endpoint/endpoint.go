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
	return CreatePostResp{
		ID:          post.ID(),
		Text:        post.Text(),
		FilesPaths:  post.FilesPaths(),
		ImagesPaths: post.ImagesPaths(),
	}, nil
}
