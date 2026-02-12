package service

import (
	"encoding/json"
	"os"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/clients"
	"social-media/internal/common/files"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/domain/message"
	"social-media/internal/posts/domain/post"
	"social-media/internal/posts/models"
	"social-media/internal/posts/repository"

	"github.com/pkg/errors"
)

type Service interface {
	CreatePost(CreatePostReq) (*post.Post, error)
	GetPosts(string, int) ([]*post.Post, error)
	UpdatePost(UpdatePostReq) (*post.Post, error)
	DeletePost(string, string) error
	MissedPostsNumber(string) ([]*post.MissedNumber, error)
	GetPostFile(string, string) (*os.File, error)

	commentsService
	chatMessagesService
}

type postsService struct {
	repo  repository.Repository
	auth  clients.AuthClient
	users clients.UsersClient
	chats clients.ChatsClient
	ai    clients.AIClient[*comment.Comment]
}

func New(r repository.Repository, auth clients.AuthClient, users clients.UsersClient,
	chat clients.ChatsClient, ai clients.AIClient[*comment.Comment]) Service {
	return &postsService{
		repo:  r,
		auth:  auth,
		users: users,
		chats: chat,
		ai:    ai,
	}
}

func (s *postsService) CreatePost(req CreatePostReq) (*post.Post, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}

	imagesPaths, err := files.Process("posts", req.Images)
	if err != nil {
		return nil, err
	}
	filesPaths, err := files.Process("posts", req.Files)
	if err != nil {
		return nil, err
	}

	post, err := s.repo.CreatePost(repository.CreatePostReq{
		CreateMessageReq: repository.CreateMessageReq{
			UserID:      id,
			Text:        req.Text,
			FilesPaths:  filesPaths,
			ImagesPaths: imagesPaths,
		},
	})
	if err != nil {
		return nil, err
	}

	msg := newPostCreatedNotification(models.PostToModel(post))
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if err = s.users.SendPostNotification(post.UserID(), string(jsonMsg)); err != nil {
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

	if p.UserID() != id {
		return nil, common.ErrForbidden
	}

	p.Update(req.Text, req.Images, req.Files, "posts", req.DeletedImages, req.DeletedFiles)

	if err := s.repo.UpdatePost(p); err != nil {
		return nil, err
	}

	notification := newPostUpdatedNotification(models.PostToModel(p))
	bytes, err := json.Marshal(notification)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	s.users.SendPostNotification(p.UserID(), string(bytes))
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

	if err := s.repo.DeletePost(id); err != nil {
		return err
	}

	notification := newPostDeletedNotification(models.PostToModel(post))
	bytes, err := json.Marshal(notification)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	s.users.SendPostNotification(post.UserID(), string(bytes))
	return nil
}

func (s *postsService) MissedPostsNumber(token string) ([]*post.MissedNumber, error) {
	userID, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	posts, err := s.users.LastReadPosts(userID)
	if err != nil {
		return nil, err
	}

	convertedPosts := slice.MustConvert(posts, func(p clients.Post) *post.LastRead {
		log.Infof("Post { ID:%s, FollowingID:%d }\n", p.ID, p.FollowingID)
		return post.NewLastReadPost(p.FollowingID, p.ID)
	})
	return s.repo.MissedPostsNumber(convertedPosts)
}

func addUsersFullNames[T message.MessageI](s *postsService, messages []T) error {
	if messages == nil {
		return nil
	}

	ids := make([]int, 0, len(messages))
	for _, comment := range messages {
		ids = append(ids, comment.UserID())
	}

	fullNames, err := s.users.UsersInfo(ids)
	if err != nil {
		return err
	}

	for i, fullName := range fullNames {
		messages[i].AddUserFullName(fullName.Name, fullName.Surname)
	}

	return nil
}

func (s *postsService) GetPostFile(token, fileID string) (*os.File, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return files.Get("posts", fileID)
}
