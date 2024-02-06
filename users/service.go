package users

import (
	"log"
	"social-media/common"
	"social-media/hash"

)

type Service interface {
	CreateUser(string, string, string, string) (string, error)
	Login(string, string) (string, error)
	GetUser(string) (*User, error)
	UpdateUser(string, string, string, string, string) error
	GetLoginsByInfo(string) ([]string, error)
	FollowUser(string, string) error
	GetFollowedLogins(string) ([]string, error)
	Log(error)
}

type userService struct {
	repo   Repository
	logger *log.Logger
	auth   common.AuthClient
}

func NewService(r Repository, logger *log.Logger, auth common.AuthClient) Service {
	return &userService{
		repo:   r,
		logger: logger,
		auth:   auth,
	}
}

func (s *userService) CreateUser(login, firstName, secondName, password string) (string, error) {
	exists, err := s.repo.UserExists(login)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}
	if exists {
		return "", common.ErrLoginExists
	}

	hashPsw, err := hash.HashPassword(password)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}

	user := NewUser(login, WithFirstName(firstName), WithSecondName(secondName), WithPassword(hashPsw))

	id, err := s.repo.SaveUser(user)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}

	user.ID = id
	user.Register()

	token, err := s.auth.GenerateToken(user.ID, user.Login)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}
	return token, nil
}

func (s *userService) Login(login, password string) (string, error) {
	exists, err := s.repo.UserExists(login)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}
	if !exists {
		return "", common.ErrWrongCredentials
	}

	id, encodedPassw, err := s.repo.GetCredentials(login)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}

	if !hash.CheckPassword(password, encodedPassw) {
		return "", common.ErrWrongCredentials
	}

	user := NewUser(login, WithID(id))
	user.Register()

	token, err := s.auth.GenerateToken(user.ID, user.Login)
	if err != nil {
		s.Log(err)
		return "", common.ErrInternal
	}
	return token, nil
}

func (s *userService) GetUser(login string) (*User, error) {
	exists, err := s.repo.UserExists(login)
	if err != nil {
		s.Log(err)
		return nil, common.ErrInternal
	}
	if !exists {
		return nil, common.ErrNotFound
	}

	user, err := s.repo.GetUser(login)
	if err != nil {
		s.Log(err)
		return nil, common.ErrInternal
	}
	return user, nil
}

func (s *userService) UpdateUser(token, firstName, secondName, bio, interests string) error {
	id, login, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	user := NewUser(login, WithID(id), WithFirstName(firstName), WithSecondName(secondName),
		WithBio(bio), WithInterests(interests))

	err = s.repo.UpdateUser(user)
	if err != nil {
		s.Log(err)
		return common.ErrInternal
	}
	return nil
}

func (s *userService) GetLoginsByInfo(info string) ([]string, error) {
	logins, err := s.repo.GetLoginsByInfo(info)
	if err != nil {
		s.Log(err)
		return nil, common.ErrInternal
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
		s.Log(err)
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
	if err != nil{
		s.Log(err)
		return nil, common.ErrInternal
	}
	return logins, nil
}

func (s *userService) Log(err error) {
	s.logger.Println(err)
}
