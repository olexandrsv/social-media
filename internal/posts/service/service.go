package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/files"
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
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}
	post, err := s.repo.CreatePost(repository.CreatePostReq{
		UserID:     id,
		Text:       text,
		FilesPaths:  filesPaths,
		ImagesPaths: imagesPaths,
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postsService) GetPosts(token string, userID int) ([]*post.Post, error) {
	_, _, err := s.auth.ValidateToken(token)
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

	addedImages, err := files.Process(req.Images)
	if err != nil {
		return nil, err
	}
	addedFiles, err := files.Process(req.Files)
	if err != nil {
		return nil, err
	}

	remainedImages, err := files.Remained(p.ImagesPaths(), req.DeletedImages)
	if err != nil {
		return nil, err
	}
	remainedFiles, err := files.Remained(p.FilesPaths(), req.DeletedFiles)
	if err != nil {
		return nil, err
	}

	if err := files.Delete(req.DeletedImages); err != nil {
		return nil, err
	}
	if err := files.Delete(req.DeletedFiles); err != nil {
		return nil, err
	}

	imagesPaths := append(remainedImages, addedImages...)
	filesPaths := append(remainedFiles, addedFiles...)

	newPost := post.New(req.ID, id, post.WithText(req.Text), post.WithImagesPaths(imagesPaths),
		post.WithFilesPaths(filesPaths))
	if err := s.repo.UpdatePost(newPost); err != nil {
		return nil, err
	}
	return newPost, nil
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

func (s *postsService) PostComments(token, id string) ([]*comment.Comment, error){
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.repo.PostComments(id)
}
