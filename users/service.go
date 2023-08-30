package users

import (
	"errors"
	"fmt"
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
}

type userService struct {
	repo Repository
	auth common.AuthClient
}

func NewService(r Repository, auth common.AuthClient) Service {
	return &userService{
		repo: r,
		auth: auth,
	}
}

func (s *userService) CreateUser(login, firstName, secondName, password string) (string, error) {
	hashPsw, err := hash.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := NewUser(login, WithFirstName(firstName), WithSecondName(secondName), WithPassword(hashPsw))

	id, err := s.repo.SaveUser(user)
	if err != nil {
		return "", err
	}

	user.ID = id
	user.Register()

	token, err := s.auth.GenerateToken(user.ID, user.Login)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) Login(login, password string) (string, error) {
	id, encodedPassw, err := s.repo.GetCredentials(login)
	if err != nil {
		return "", err
	}
	
	if !hash.CheckPassword(password, encodedPassw) {
		return "", errors.New("forbidden")
	}

	user := NewUser(login, WithID(id))
	user.Register()

	token, err := s.auth.GenerateToken(user.ID, user.Login)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetUser(login string) (*User, error){
	user, err := s.repo.GetUser(login)
	if err != nil{
		return nil, err
	}
	fmt.Println(user)
	return user, nil
}

func (s *userService) UpdateUser(token, firstName, secondName, bio, interests string) error{
	id, login, err := s.auth.ValidateToken(token)
	if err != nil{
		return err
	}
	user := NewUser(login, WithID(id), WithFirstName(firstName), WithSecondName(secondName), 
		WithBio(bio), WithInterests(interests))

	err = s.repo.UpdateUser(user)
	return err
}

func (s *userService) GetLoginsByInfo(info string) ([]string, error){
	return s.repo.GetLoginsByInfo(info)
}

func (s *userService) FollowUser(token, followedLogin string) error{
	id, _, err := s.auth.ValidateToken(token)
	if err != nil{
		return err
	}
	followedID, err := s.repo.GetIdByLogin(followedLogin)
	if err != nil{
		return err
	}
	return s.repo.Subscribe(id, followedID)
}

func (s *userService) GetFollowedLogins(token string) ([]string, error){
	id, _, err := s.auth.ValidateToken(token)
	if err != nil{
		return nil, err
	}
	return s.repo.GetFollowedLogins(id)
}