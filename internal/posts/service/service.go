package service

import (
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/clients"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/domain/post"
	"social-media/internal/posts/repository"

	"github.com/pkg/errors"
)

type Service interface {
	CreatePost(token, text string, filesPaths, imagesPaths []string) (*post.Post, error)
	GetPosts(string, int) ([]*post.Post, error)
	UpdatePost(UpdatePostReq) (*post.Post, error)
	DeletePost(string, string) error
	commentsService
}

type postsService struct {
	repo repository.Repository
	auth clients.AuthClient
	users clients.UsersClient
}

func New(r repository.Repository, auth clients.AuthClient, users clients.UsersClient) Service {
	return &postsService{
		repo: r,
		auth: auth,
		users: users,
	}
}

func (s *postsService) CreatePost(token, text string, filesPaths, imagesPaths []string) (*post.Post, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}
	post, err := s.repo.CreatePost(repository.CreatePostReq{
		CreateMessageReq: repository.CreateMessageReq{
			UserID:      id,
			Text:        text,
			FilesPaths:  filesPaths,
			ImagesPaths: imagesPaths,
		},
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postsService) GetPosts(token string, userID int) ([]*post.Post, error) {
	fullNames, err := s.users.UsersFullNames([]int{1, 2, 3})
	if err != nil {
		return nil, err
	}
	log.Error(fmt.Errorf("%#v", fullNames))
	
	_, _, err = s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.repo.UserPosts(userID)
}

func (s *postsService) UpdatePost(req UpdatePostReq) (*post.Post, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.GetPost(req.ID)
	if err != nil {
		return nil, err
	}

	if p.UserID() != id {
		return nil, common.ErrForbidden
	}

	p.Update(req.Text, req.Images, req.Files, req.DeletedImages, req.DeletedFiles)

	if err := s.repo.UpdatePost(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *postsService) DeletePost(token, id string) error {
	userID, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	post, err := s.repo.GetPost(id)
	if err != nil {
		return err
	}

	if post.UserID() != userID {
		return common.ErrForbidden
	}

	return s.repo.DeletePost(id)
}

func (s *postsService) PostComments(token, id string) ([]*comment.Comment, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.repo.PostComments(id)
}
