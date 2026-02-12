package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/clients"
	"social-media/internal/common/connection"
	"social-media/internal/common/slice"
	"social-media/internal/users/domain/post"
	"social-media/internal/users/domain/user"
	"social-media/internal/users/repository"

	"github.com/pkg/errors"
)

type Service interface {
	CreateUser(string, string, string, string) (string, int, error)
	Login(string, string) (string, int, error)
	GetUser(string, int) (*user.User, error)
	UpdateUser(string, string, string, string, string) error
	GetUsersByInfo(string) ([]*user.User, error)
	FollowUser(string, int) error
	GetFollowedUsers(string) ([]*user.User, error)
	UsersInfo([]int) ([]*user.User, error)
	UpdateReadPosts(string, int, string) error
	LastReadPosts(int) ([]*post.Post, error)
	SendPostNotification(int, string) error
	SetUserConnection(string, connection.MessageConnection) error
	SendNotification(userID int, receiversIDs []int, message string) error
}

type userService struct {
	repo repository.Repository
	auth clients.AuthClient
}

func New(r repository.Repository, auth clients.AuthClient) Service {
	return &userService{
		repo: r,
		auth: auth,
	}
}

func (s *userService) CreateUser(login, name, surname, password string) (string, int, error) {
	var userID int
	var token string
	exists, err := s.repo.UserExists(login)
	if err != nil {
		return token, userID, err
	}
	if exists {
		return token, userID, common.ErrLoginExists
	}

	hashPsw, err := hashPassword(password)
	if err != nil {
		return token, userID, err
	}

	userModel := repository.NewUserModel(login, name, surname, hashPsw)

	u, err := s.repo.CreateUser(userModel)
	if err != nil {
		return token, userID, err
	}

	token, err = s.auth.GenerateToken(u.ID(), u.Login())
	if err != nil {
		return token, userID, err
	}

	userID = u.ID()
	user.Set(u)
	return token, userID, nil
}

func (s *userService) Login(login, password string) (string, int, error) {
	id, encodedPassw, err := s.repo.GetCredentials(login)
	if err != nil {
		return "", 0, err
	}

	if !checkPassword(password, encodedPassw) {
		return "", 0, common.ErrWrongCredentials
	}

	u := user.New(id, login)

	token, err := s.auth.GenerateToken(u.ID(), u.Login())
	if err != nil {
		return "", 0, err
	}

	user.Set(u)
	return token, u.ID(), nil
}

func (s *userService) GetUser(token string, id int) (*user.User, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUser(id)
}

func (s *userService) UpdateUser(token, name, surname, bio, interests string) error {
	id, login, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	user := user.New(id, login, user.WithName(name), user.WithSurname(surname),
		user.WithBio(bio), user.WithInterests(interests))

	err = s.repo.UpdateUser(user)
	if err != nil {
		log.Error(err)
		return common.ErrInternal
	}
	return nil
}

func (s *userService) GetUsersByInfo(info string) ([]*user.User, error) {
	return s.repo.GetUsersByInfo(info)
}

func (s *userService) FollowUser(token string, followedID int) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	if id == followedID {
		return common.ErrInvalidData
	}

	exists, err := s.repo.SubscriptionExists(id, followedID)
	if err != nil {
		return err
	}
	if exists {
		return common.ErrSubscriptionExists
	}
	return s.repo.Subscribe(id, followedID)
}

func (s *userService) GetFollowedUsers(token string) ([]*user.User, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.repo.GetFollowedUsers(id)
}

func (s *userService) UsersInfo(ids []int) ([]*user.User, error) {
	return s.repo.UsersInfo(ids)
}

func (s *userService) UpdateReadPosts(token string, ownerID int, lastReadPostID string) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	return s.repo.UpdateReadPosts(ownerID, id, lastReadPostID)
}

func (s *userService) LastReadPosts(userID int) ([]*post.Post, error) {
	return s.repo.LastReadPosts(userID)
}

func (s *userService) SendPostNotification(userID int, message string) error {
	followers, err := s.repo.GetFollowers(userID)
	if err != nil {
		return err
	}
	followersIDs := slice.MustConvert(followers, func(u *user.User) int {
		return u.ID()
	})
	return s.SendNotification(userID, followersIDs, message)
}

func (s *userService) SendNotification(userID int, receiversIDs []int, message string) error {
	for _, receiver := range receiversIDs {
		u, ok := user.Get(receiver)
		if !ok {
			continue
		}
		if err := u.Send(message); err != nil {
			log.Error(errors.WithStack(err))
			return common.ErrInternal
		}
	}
	return nil
}

func (s *userService) SetUserConnection(token string, conn connection.MessageConnection) error {
	id, login, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	u, exist := user.Get(id)
	if !exist {
		u = user.New(id, login)
		user.Set(u)
	}

	u.SetConnection(conn)
	return nil
}
