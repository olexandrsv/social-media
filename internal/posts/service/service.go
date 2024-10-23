package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/domain/post"
	"social-media/internal/posts/repository"
)

type Service interface {
	CreatePost(token, text string, filesPaths, imagesPaths []string) (*post.Post, error)
}

type postsService struct {
	repo repository.Repository
	auth common.AuthClient
}

func New(r repository.Repository, auth common.AuthClient) Service {
	return &postsService{
		repo: r,
		auth: auth,
	}
}

func (s *postsService) CreatePost(token, text string, filesPaths, imagesPaths []string) (*post.Post, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		log.Error(err)
		return nil, common.ErrInvalidToken
	}
	post, err := s.repo.CreatePost(repository.PostModel{
		UserID:     id,
		Text:       text,
		FilesPath:  filesPaths,
		ImagesPath: imagesPaths,
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}
