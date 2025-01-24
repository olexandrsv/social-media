package service

import (
	"social-media/internal/common"
	"social-media/internal/common/files"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/repository"
)

type commentsService interface {
	PostComments(string, string) ([]*comment.Comment, error)
	CreatePostComment(CreatePostCommentReq) (*comment.Comment, error)
	DeletePostComment(string, string, string) error

	UpdateComment(UpdateCommentReq) (*comment.Comment, error)

	CommentComments(string, string) ([]*comment.Comment, error)
	CreateCommentComment(CreateCommentCommentReq) (*comment.Comment, error)
	DeleteCommentComment(string, string, string) error
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

	comment, err := s.repo.CreatePostComment(repository.CreatePostCommentReq{
		PostID: req.PostID,
		CreateMessageReq: repository.CreateMessageReq{
			UserID:      id,
			Text:        req.Text,
			ImagesPaths: imagesPaths,
			FilesPaths:  filesPaths,
		},
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *postsService) DeletePostComment(token, parentID, commentID string) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	comment, err := s.repo.GetComment(commentID)
	if err != nil {
		return err
	}

	if comment.UserID() != id {
		return common.ErrForbidden
	}

	return s.repo.DeletePostComment(parentID, commentID)
}

func (s *postsService) UpdateComment(req UpdateCommentReq) (*comment.Comment, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	c, err := s.repo.GetComment(req.ID)
	if err != nil {
		return nil, err
	}

	if c.UserID() != id {
		return nil, common.ErrForbidden
	}

	addedImages, err := files.Process(req.Images)
	if err != nil {
		return nil, err
	}
	addedFiles, err := files.Process(req.Files)
	if err != nil {
		return nil, err
	}

	remainedImages, err := files.Remained(c.ImagesPaths(), req.DeletedImages)
	if err != nil {
		return nil, err
	}
	remainedFiles, err := files.Remained(c.FilesPaths(), req.DeletedFiles)
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

	s.repo.UpdateComment()

	return comment.New(req.ID, id, comment.WithText(req.Text), comment.WithImagesPaths(imagesPaths),
		comment.WithFilesPaths(filesPaths)), nil
}

func (s *postsService) CommentComments(token, commentID string) ([]*comment.Comment, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.repo.CommentComments(commentID)
}

func (s *postsService) CreateCommentComment(req CreateCommentCommentReq) (*comment.Comment, error) {
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

	return s.repo.CreateCommentComment(repository.CreateCommentCommentReq{
		CommentID: req.ParentID,
		CreateMessageReq: repository.CreateMessageReq{
			UserID:      id,
			Text:        req.Text,
			ImagesPaths: imagesPaths,
			FilesPaths:  filesPaths,
		},
	})
}

func (s *postsService) DeleteCommentComment(token, parentID, commentID string) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	comment, err := s.repo.GetComment(commentID)
	if err != nil {
		return err
	}

	if comment.UserID() != id {
		return common.ErrForbidden
	}

	return s.repo.DeleteCommentComment(parentID, commentID)
}