package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/clients"
	"social-media/internal/users/domain/user"
	"social-media/internal/users/repository"
)

type Service interface {
	CreateUser(string, string, string, string) (string, int, error)
	Login(string, string) (string, int, error)
	GetUser(string, int) (*user.User, error)
	UpdateUser(string, string, string, string, string) error
	GetUsersByInfo(string) ([]*user.User, error)
	FollowUser(string, int) error
	GetFollowedUsers(string) ([]*user.User, error)
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
	exists, err := s.repo.UserExists(login)
	if err != nil {
		return "", 0, err
	}
	if exists {
		return "", 0, common.ErrLoginExists
	}

	hashPsw, err := hashPassword(password)
	if err != nil {
		return "", 0, err
	}

	userModel := repository.NewUserModel(login, name, surname, hashPsw)

	user, err := s.repo.CreateUser(userModel)
	if err != nil {
		return "", 0, err
	}

	user.Register()

	token, err := s.auth.GenerateToken(user.ID(), user.Login())
	if err != nil {
		return "", 0, err
	}
	return token, user.ID(), nil
}

func (s *userService) Login(login, password string) (string, int, error) {
	id, encodedPassw, err := s.repo.GetCredentials(login)
	if err != nil {
		return "", 0, err
	}

	if !checkPassword(password, encodedPassw) {
		return "", 0, common.ErrWrongCredentials
	}

	user := user.New(id, login)
	user.Register()

	token, err := s.auth.GenerateToken(user.ID(), user.Login())
	if err != nil {
		return "", 0, err
	}
	return token, user.ID(), nil
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
