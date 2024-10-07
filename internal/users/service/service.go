package service

import (
	"database/sql"
	"errors"
	"social-media/hash"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/users/domain/user"
	"social-media/internal/users/repository"
)



type Service interface {
	CreateUser(string, string, string, string) (string, error)
	Login(string, string) (string, error)
	GetUser(string) (*user.User, error)
	UpdateUser(string, string, string, string, string) error
	GetLoginsByInfo(string) ([]string, error)
	FollowUser(string, string) error
	GetFollowedLogins(string) ([]string, error)
}

type userService struct {
	repo   repository.Repository
	auth   common.AuthClient
}

func New(r repository.Repository, auth common.AuthClient) Service {
	return &userService{
		repo:   r,
		auth:   auth,
	}
}

func (s *userService) CreateUser(login, name, surname, password string) (string, error) {
	exists, err := s.repo.UserExists(login)
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}
	if exists {
		return "", common.ErrLoginExists
	}

	hashPsw, err := hash.HashPassword(password)
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}

	userModel := repository.NewUserModel(login, name, surname, hashPsw)

	user, err := s.repo.CreateUser(userModel)
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}

	user.Register()

	token, err := s.auth.GenerateToken(user.ID(), user.Login())
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}
	return token, nil
}

func (s *userService) Login(login, password string) (string, error) {
	id, encodedPassw, err := s.repo.GetCredentials(login)
	if errors.Is(err, sql.ErrNoRows) {
		log.Error(err)
		return "", common.ErrWrongCredentials
	}
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}

	if !hash.CheckPassword(password, encodedPassw) {
		return "", common.ErrWrongCredentials
	}

	user := user.NewUser(id, login)
	user.Register()

	token, err := s.auth.GenerateToken(user.ID(), user.Login())
	if err != nil {
		log.Error(err)
		return "", common.ErrInternal
	}
	return token, nil
}

func (s *userService) GetUser(login string) (*user.User, error) {
	user, err := s.repo.GetUser(login)
	if errors.Is(err, sql.ErrNoRows) {
		log.Error(err)
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(err)
		return nil, common.ErrInternal
	}
	return user, nil
}

func (s *userService) UpdateUser(token, name, surname, bio, interests string) error {
	id, login, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	user := user.NewUser(id, login, user.WithName(name), user.WithSurname(surname),
		user.WithBio(bio), user.WithInterests(interests))

	err = s.repo.UpdateUser(user)
	if err != nil {
		log.Error(err)
		return common.ErrInternal
	}
	return nil
}

func (s *userService) GetLoginsByInfo(info string) ([]string, error) {
	logins, err := s.repo.GetLoginsByInfo(info)
	if err != nil {
		log.Error(err)
		return nil, common.ErrInternal
	}
	if logins == nil {
		return nil, common.ErrNotFound
	}
	return logins, nil
}

func (s *userService) FollowUser(token, followedLogin string) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	err = s.repo.Subscribe(id, followedLogin)
	if err != nil {
		log.Error(err)
		return common.ErrInternal
	}
	return nil
}

func (s *userService) GetFollowedLogins(token string) ([]string, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	logins, err := s.repo.GetFollowedLogins(id)
	if err != nil {
		log.Error(err)
		return nil, common.ErrInternal
	}
	return logins, nil
}

