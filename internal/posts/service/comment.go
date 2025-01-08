package service

import (
	"social-media/internal/common/files"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/repository"
)

type commentsService interface {
	PostComments(string, string) ([]*comment.Comment, error)
	CreatePostComment(CreatePostCommentReq) (*comment.Comment, error)
}

func (s *postsService) CreatePostComment(req CreatePostCommentReq) (*comment.Comment, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	imagesPaths, err := files.Process(req.Images)
	if err != nil {
		return nil, err
	}
	filesPaths, err := files.Process(req.Files)
	if err != nil {
		return nil, err
	}
	
	comment, err := s.repo.CreatePostComment(repository.CreateCommentReq{
		UserID:     id,
		PostID:     req.PostID,
		Text:       req.Text,
		ImagesPath: imagesPaths,
		FilesPath:  filesPaths,
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

