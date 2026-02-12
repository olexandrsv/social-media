package service

import (
	"social-media/internal/common"
	"social-media/internal/common/files"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/domain/tone"
	"social-media/internal/posts/repository"
)

type commentsService interface {
	PostComments(string, string) (*tone.Tone, []*comment.Comment, error)
	CreatePostComment(CreatePostCommentReq) (*comment.Comment, error)
	DeletePostComment(string, string, string) error

	UpdateComment(UpdateCommentReq) (*comment.Comment, error)

	CommentComments(string, string) ([]*comment.Comment, error)
	CreateCommentComment(CreateCommentCommentReq) (*comment.Comment, error)
	DeleteCommentComment(string, string, string) error
}

func (s *postsService) PostComments(token, id string) (*tone.Tone, []*comment.Comment, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, nil, err
	}

	comments, err := s.repo.PostComments(id)
	if err != nil {
		return nil, nil, err
	}

	if err := addUsersFullNames(s, comments); err != nil {
		return nil, nil, err
	}

	tone, err := s.ai.EstimateTone(comments)
	if err != nil {
		return nil, nil, err
	}

	return tone, comments, nil
}

func (s *postsService) CreatePostComment(req CreatePostCommentReq) (*comment.Comment, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	imagesPaths, err := files.Process("posts", req.Images)
	if err != nil {
		return nil, err
	}
	filesPaths, err := files.Process("posts", req.Files)
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

	comment, err := s.repo.GetComment(req.ID)
	if err != nil {
		return nil, err
	}

	if comment.UserID() != id {
		return nil, common.ErrForbidden
	}

	err = comment.Update(req.Text, req.Images, req.Files, "posts", req.DeletedImages, req.DeletedFiles)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *postsService) CommentComments(token, commentID string) ([]*comment.Comment, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	comments, err := s.repo.CommentComments(commentID)
	if err != nil {
		return nil, err
	}

	addUsersFullNames(s, comments)

	return comments, nil
}

func (s *postsService) CreateCommentComment(req CreateCommentCommentReq) (*comment.Comment, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	imagesPaths, err := files.Process("posts", req.Images)
	if err != nil {
		return nil, err
	}
	filesPaths, err := files.Process("posts", req.Files)
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
