package endpoint

import (
	"context"
	"errors"
	"social-media/internal/common"
	"social-media/internal/users/domain/user"
	"testing"
)

type mockService struct {
	createUser        func(string, string, string, string) (string, error)
	login             func(string, string) (string, error)
	getUser           func(string) (*user.User, error)
	updateUser        func(string, string, string, string, string) error
	getLoginsByInfo   func(string) ([]string, error)
	followUser        func(string, string) error
	getFollowedLogins func(string) ([]string, error)
}

func (s *mockService) CreateUser(login, firstName, secondName, password string) (string, error) {
	return s.createUser(login, firstName, secondName, password)
}

func (s *mockService) Login(login, password string) (string, error) {
	return s.login(login, password)
}

func (s *mockService) GetUser(login string) (*user.User, error) {
	return s.getUser(login)
}

func (s *mockService) UpdateUser(token, firstName, secondName, bio, interests string) error {
	return s.updateUser(token, firstName, secondName, bio, interests)
}

func (s *mockService) GetLoginsByInfo(info string) ([]string, error) {
	return s.getLoginsByInfo(info)
}

func (s *mockService) FollowUser(token, followedLogin string) error {
	return s.followUser(token, followedLogin)
}

func (s *mockService) GetFollowedLogins(token string) ([]string, error) {
	return s.getFollowedLogins(token)
}

func (s *mockService) Log(err error) {

}

func TestCreateUser(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), 
		user.WithSurname("Smith"), user.WithPassword("4grts67Jk"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "")
	oliver := user.NewUser(4, "oliver", user.WithName("Oliver"), 
		user.WithSurname("Brown"), user.WithPassword("j;gsfhgtuw432"))

	checkParams := func(login, firstName, secondName, password string, u *user.User) {
		if login != u.Login() || firstName != u.Name() || secondName != u.Surname() || password != u.Password() {
			t.Errorf("error login %q, firstName %q, secondName %q, password %q; epxected: %q, %q, %q, %q",
				login, firstName, secondName, password, u.Login(), u.Name(), u.Surname(), u.Password())
		}
	}
	e := errors.New("endpoint error")

	data := []struct {
		s      *mockService
		req    interface{}
		resp   interface{}
		resErr error
	}{
		{
			s: &mockService{
				createUser: func(login, firstName, secondName, password string) (string, error) {
					checkParams(login, firstName, secondName, password, bob)
					return firstName + secondName, nil
				},
			},
			req: CreateUserReq{
				Login:      bob.Login(),
				Name:  bob.Name(),
				Surname: bob.Surname(),
				Password:   bob.Password(),
			},
			resp: AuthResp{
				Token: bob.Name() + bob.Surname(),
			},
			resErr: nil,
		},
		{
			req:    nil,
			resp:   nil,
			resErr: common.ErrInternal,
		},
		{
			s: &mockService{
				createUser: func(login, firstName, secondName, password string) (string, error) {
					checkParams(login, firstName, secondName, password, ben)
					return firstName+secondName, nil
				},
			},
			req: CreateUserReq{
				Login:      ben.Login(),
				Name:  ben.Name(),
				Surname: ben.Surname(),
				Password:   ben.Password(),
			},
			resp: nil,
			resErr: errors.New("Field password can't be less than 8 symbols"),
		},
		{
			s: &mockService{
				createUser: func(login, firstName, secondName, password string) (string, error) {
					checkParams(login, firstName, secondName, password, bill)
					return firstName + secondName, nil
				},
			},
			req: CreateUserReq{
				Login:      bill.Login(),
				Name:  bill.Name(),
				Surname: bill.Surname(),
				Password:   bill.Password(),
			},
			resp:   nil,
			resErr: errors.New("Field login can't be empty\nField password can't be less than 8 symbols"),
		},
		{
			s: &mockService{
				createUser: func(login, firstName, secondName, password string) (string, error) {
					checkParams(login, firstName, secondName, password, oliver)
					return "", e
				},
			},
			req: CreateUserReq{
				Login:      oliver.Login(),
				Name:  oliver.Name(),
				Surname: oliver.Surname(),
				Password:   oliver.Password(),
			},
			resp:   nil,
			resErr: e,
		},
	}

	for _, d := range data {
		e := NewEndpoints(d.s)
		resp, err := e.CreateUser(context.Background(), d.req)

		if err != d.resErr && err.Error() != d.resErr.Error() {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if resp != d.resp {
			t.Errorf("wrong resp %+v, expected: %+v", resp, d.resp)
		}
	}
}


