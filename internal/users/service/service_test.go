package service

import (
	"database/sql"
	"errors"
	"reflect"
	"social-media/hash"
	"social-media/internal/common"
	"social-media/internal/common/app"
	"social-media/internal/users/domain/user"
	"social-media/internal/users/repository"
	"testing"
)

type mockRepo struct {
	createUser        func(repository.UserModel) (*user.User, error)
	getCredentials    func(string) (int, string, error)
	getUser           func(string) (*user.User, error)
	updateUser        func(*user.User) error
	userExists        func(string) (bool, error)
	getLoginsByInfo   func(string) ([]string, error)
	subscribe         func(int, string) error
	getFollowedLogins func(int) ([]string, error)
}

func (r *mockRepo) CreateUser(u repository.UserModel) (*user.User, error) {
	return r.createUser(u)
}

func (r *mockRepo) GetCredentials(login string) (int, string, error) {
	return r.getCredentials(login)
}

func (r *mockRepo) GetUser(login string) (*user.User, error) {
	return r.getUser(login)
}

func (r *mockRepo) UpdateUser(u *user.User) error {
	return r.updateUser(u)
}

func (r *mockRepo) UserExists(login string) (bool, error) {
	return r.userExists(login)
}

func (r *mockRepo) GetLoginsByInfo(info string) ([]string, error) {
	return r.getLoginsByInfo(info)
}

func (r *mockRepo) Subscribe(userID int, followedLogin string) error {
	return r.subscribe(userID, followedLogin)
}

func (r *mockRepo) GetFollowedLogins(id int) ([]string, error) {
	return r.getFollowedLogins(id)
}

type mockAuthClient struct {
	generateToken func(int, string) (string, error)
	validateToken func(string) (int, string, error)
}

func (c *mockAuthClient) GenerateToken(id int, login string) (string, error) {
	return c.generateToken(id, login)
}

func (c *mockAuthClient) ValidateToken(token string) (int, string, error) {
	return c.validateToken(token)
}


func TestCreateUser(t *testing.T) {
	app.InitUsersService()
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))
	oliver := user.NewUser(4, "oliver", user.WithName("Oliver"), user.WithSurname("Brown"))
	henry := user.NewUser(5, "henry", user.WithName("Henry"), user.WithSurname("Williams"))

	data := []struct {
		repo     *mockRepo
		client   *mockAuthClient
		user     *user.User
		resToken string
		resErr   error
	}{
		{
			repo: &mockRepo{
				userExists: func(login string) (bool, error) {
					if login != bob.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bob.Login())
					}
					return false, nil
				},
				createUser: func(u repository.UserModel) (*user.User, error) {
					if u.Login != bob.Login() || u.Name != bob.Name() || u.Surname != bob.Surname() {
						t.Errorf("wrong user %+v, expected: %+v", u, bob)
					}
					return user.NewUser(bob.ID(), u.Login, user.WithName(u.Name),
						user.WithSurname(u.Surname), user.WithPassword(u.Password)), nil
				},
			},
			client: &mockAuthClient{
				generateToken: func(id int, login string) (string, error) {
					if id != bob.ID() || login != bob.Login() {
						t.Errorf("wrong id %q or login %q expected: %q, %q", id, login, bob.ID(), bob.Login())
					}
					return "token", nil
				},
			},
			user:     bob,
			resToken: "token",
			resErr:   nil,
		},
		{
			repo: &mockRepo{
				userExists: func(login string) (bool, error) {
					if login != ben.Login() {
						t.Errorf("wrong login %q, expected: %q", login, ben.Login())
					}
					return false, errors.New("service error")
				},
			},
			user:     ben,
			resToken: "",
			resErr:   common.ErrInternal,
		},
		{
			repo: &mockRepo{
				userExists: func(login string) (bool, error) {
					if login != bill.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bill.Login())
					}
					return true, nil
				},
			},
			user:     bill,
			resToken: "",
			resErr:   common.ErrLoginExists,
		},
		{
			repo: &mockRepo{
				userExists: func(login string) (bool, error) {
					if login != oliver.Login() {
						t.Errorf("wrong login %q, expected: %q", login, oliver.Login())
					}
					return false, nil
				},
				createUser: func(u repository.UserModel) (*user.User, error) {
					if u.Login != oliver.Login() || u.Name != oliver.Name() || u.Surname != oliver.Surname() {
						t.Errorf("wrong user %+v, expected: %+v", u, bob)
					}
					return nil, errors.New("service error")
				},
			},
			user:     oliver,
			resToken: "",
			resErr:   common.ErrInternal,
		},
		{
			repo: &mockRepo{
				userExists: func(login string) (bool, error) {
					if login != henry.Login() {
						t.Errorf("wrong login %q, expected: %q", login, henry.Login())
					}
					return false, nil
				},
				createUser: func(u repository.UserModel) (*user.User, error) {
					if u.Login != henry.Login() || u.Name != henry.Name() || u.Surname != henry.Surname() {
						t.Errorf("wrong user %+v, expected: %+v", u, henry)
					}
					return user.NewUser(henry.ID(), u.Login, user.WithName(u.Name),
						user.WithSurname(u.Surname), user.WithPassword(u.Password)), nil
				},
			},
			client: &mockAuthClient{
				generateToken: func(id int, login string) (string, error) {
					if id != henry.ID() || login != henry.Login() {
						t.Errorf("wrong id %q or login %q expected: %q, %q", id, login, henry.ID(), henry.Login())
					}
					return "", errors.New("service error")
				},
			},
			user:     henry,
			resToken: "",
			resErr:   common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		token, err := service.CreateUser(d.user.Login(), d.user.Name(), d.user.Surname(), d.user.Password())

		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if token != d.resToken {
			t.Errorf("wrong token %s, expected: %s", token, d.resToken)
		}
	}
}

func TestLogin(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"), user.WithPassword("xHtRe"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))
	oliver := user.NewUser(4, "oliver", user.WithName("Oliver"), user.WithSurname("Brown"),
		user.WithPassword("ifjrehT"))

	data := []struct {
		repo     *mockRepo
		client   *mockAuthClient
		user     *user.User
		resToken string
		resErr   error
	}{
		{
			repo: &mockRepo{
				getCredentials: func(login string) (int, string, error) {
					if login != bob.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bob.Login())
					}
					pwd, err := hash.HashPassword(bob.Password())
					if err != nil {
						t.Errorf("hash error: password %s", bob.Password())
					}
					return bob.ID(), pwd, nil
				},
			},
			client: &mockAuthClient{
				generateToken: func(id int, login string) (string, error) {
					if id != bob.ID() || login != bob.Login() {
						t.Errorf("wrong id %q or login %q expected: %q, %q", id, bob.ID(), login, bob.Login())
					}
					return "token", nil
				},
			},
			user:     bob,
			resToken: "token",
			resErr:   nil,
		},
		{
			repo: &mockRepo{
				getCredentials: func(login string) (int, string, error) {
					if login != bill.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bill.Login())
					}
					return 0, "", sql.ErrNoRows
				},
			},
			user:     bill,
			resToken: "",
			resErr:   common.ErrWrongCredentials,
		},
		{
			repo: &mockRepo{
				getCredentials: func(login string) (int, string, error) {
					if login != oliver.Login() {
						t.Errorf("wrong login %q, expected: %q", login, oliver.Login())
					}
					return oliver.ID(), "", errors.New("service error")
				},
			},
			user:     oliver,
			resToken: "",
			resErr:   common.ErrInternal,
		},
		{
			repo: &mockRepo{
				getCredentials: func(login string) (int, string, error) {
					if login != bob.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bob.Login())
					}
					pwd, err := hash.HashPassword(bob.Password())
					if err != nil {
						t.Errorf("hash error: password %s", bob.Password())
					}
					return bob.ID(), pwd + "invalid data", nil
				},
			},
			client: &mockAuthClient{
				generateToken: func(id int, login string) (string, error) {
					if id != bob.ID() || login != bob.Login() {
						t.Errorf("wrong id %q or login %q expected: %q, %q", id, bob.ID(), login, bob.Login())
					}
					return "token", nil
				},
			},
			user:     bob,
			resToken: "",
			resErr:   common.ErrWrongCredentials,
		},
		{
			repo: &mockRepo{
				getCredentials: func(login string) (int, string, error) {
					if login != oliver.Login() {
						t.Errorf("wrong login %q, expected: %q", login, oliver.Login())
					}
					pwd, err := hash.HashPassword(oliver.Password())
					if err != nil {
						t.Errorf("hash error: password %s", oliver.Password())
					}
					return oliver.ID(), pwd, nil
				},
			},
			client: &mockAuthClient{
				generateToken: func(id int, login string) (string, error) {
					if id != oliver.ID() || login != oliver.Login() {
						t.Errorf("wrong id %q or login %q expected: %q, %q", id, oliver.ID(), login, oliver.Login())
					}
					return "", errors.New("service error")
				},
			},
			user:     oliver,
			resToken: "",
			resErr:   common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		token, err := service.Login(d.user.Login(), d.user.Password())

		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if token != d.resToken {
			t.Errorf("wrong token %s, expected: %s", token, d.resToken)
		}
	}
}

func TestGetUser(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))

	data := []struct {
		repo    *mockRepo
		client  *mockAuthClient
		user    *user.User
		resUser *user.User
		resErr  error
	}{
		{
			repo: &mockRepo{
				getUser: func(login string) (*user.User, error) {
					if login != bob.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bob.Login())
					}
					return bob, nil
				},
			},
			user:    bob,
			resUser: bob,
			resErr:  nil,
		},
		{
			repo: &mockRepo{
				getUser: func(login string) (*user.User, error) {
					if login != ben.Login() {
						t.Errorf("wrong login %q, expected: %q", login, ben.Login())
					}
					return nil, sql.ErrNoRows
				},
			},
			user:    ben,
			resUser: nil,
			resErr:  common.ErrNotFound,
		},
		{
			repo: &mockRepo{
				getUser: func(login string) (*user.User, error) {
					if login != bill.Login() {
						t.Errorf("wrong login %q, expected: %q", login, bill.Login())
					}
					return nil, errors.New("service error")
				},
			},
			user:    bill,
			resUser: nil,
			resErr:  common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		user, err := service.GetUser(d.user.Login())

		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if user != d.resUser {
			t.Errorf("wrong user %+v, expected: %+v", user, d.resUser)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))
	e := errors.New("service error")

	data := []struct {
		repo   *mockRepo
		client *mockAuthClient
		user   *user.User
		resErr error
	}{
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					return bob.ID(), bob.Login(), nil
				},
			},
			repo: &mockRepo{
				updateUser: func(u *user.User) error {
					if !reflect.DeepEqual(u, bob) {
						t.Errorf("wrong user %+v, expected: %+v", u, bob)
					}
					return nil
				},
			},
			user:   bob,
			resErr: nil,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					return 0, "", e
				},
			},
			user:   ben,
			resErr: e,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					return bill.ID(), bill.Login(), nil
				},
			},
			repo: &mockRepo{
				updateUser: func(u *user.User) error {
					if !reflect.DeepEqual(u, bill) {
						t.Errorf("wrong user %+v, expected: %+v", u, bill)
					}
					return errors.New("service error")
				},
			},
			user:   bill,
			resErr: common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		err := service.UpdateUser(d.user.Login(), d.user.Name(), d.user.Surname(), d.user.Bio(), d.user.Interests())

		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}
	}
}

func TestGetLoginsByInfo(t *testing.T) {
	data := []struct {
		repo      *mockRepo
		client    *mockAuthClient
		info      string
		resLogins []string
		resErr    error
	}{
		{
			repo: &mockRepo{
				getLoginsByInfo: func(info string) ([]string, error) {
					return []string{"bob", "ben"}, nil
				},
			},
			info:      "book",
			resLogins: []string{"bob", "ben"},
			resErr:    nil,
		},
		{
			repo: &mockRepo{
				getLoginsByInfo: func(info string) ([]string, error) {
					return nil, errors.New("service error")
				},
			},
			info:      "student",
			resLogins: nil,
			resErr:    common.ErrInternal,
		},
		{
			repo: &mockRepo{
				getLoginsByInfo: func(info string) ([]string, error) {
					return nil, nil
				},
			},
			info:      "pupil",
			resLogins: nil,
			resErr:    common.ErrNotFound,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		logins, err := service.GetLoginsByInfo(d.info)

		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if len(logins) != len(d.resLogins) {
			t.Errorf("wrong logins %q, expected: %q", logins, d.resLogins)
			return
		}

		for i, v := range logins {
			if v != d.resLogins[i] {
				t.Errorf("wrong logins %q, expected: %q", logins, d.resLogins)
				return
			}
		}
	}
}

func TestFollowUser(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))

	data := []struct {
		repo         *mockRepo
		client       *mockAuthClient
		user         *user.User
		followedUser *user.User
		resErr       error
	}{
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != bob.Name()+bob.Surname() {
						t.Errorf("wrong token %s, expected: %s", token, bob.Name()+bob.Surname())
					}
					return bob.ID(), bob.Login(), nil
				},
			},
			repo: &mockRepo{
				subscribe: func(id int, followedLogin string) error {
					if id != bob.ID() || followedLogin != ben.Login() {
						t.Errorf("wrong id %d or login %q, expected: %d, %q", id, followedLogin, bob.ID(), bob.Login())
					}
					return nil
				},
			},
			user:         bob,
			followedUser: ben,
			resErr:       nil,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != ben.Name()+ben.Surname() {
						t.Errorf("wrong error %s, expected: %s", token, ben.Name()+ben.Surname())
					}
					return 0, "", common.ErrInvalidToken
				},
			},
			user:         ben,
			followedUser: bob,
			resErr:       common.ErrInvalidToken,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != bill.Name()+bill.Surname() {
						t.Errorf("wrong token %s, expected: %s", token, bill.Name()+bill.Surname())
					}
					return bill.ID(), bill.Login(), nil
				},
			},
			repo: &mockRepo{
				subscribe: func(id int, followedLogin string) error {
					if id != bill.ID() || followedLogin != ben.Login() {
						t.Errorf("wrong id %d or login %q, expected: %d, %q", id, followedLogin, bob.ID(), ben.Login())
					}
					return errors.New("service error")
				},
			},
			user:         bill,
			followedUser: ben,
			resErr:       common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		err := service.FollowUser(d.user.Name()+d.user.Surname(), d.followedUser.Login())
		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}
	}
}

func TestGetFollowedLogins(t *testing.T) {
	bob := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"))
	ben := user.NewUser(2, "ben", user.WithName("Ben"), user.WithSurname("Jones"))
	bill := user.NewUser(3, "bill", user.WithName("Bill"), user.WithSurname("White"))

	data := []struct {
		repo      *mockRepo
		client    *mockAuthClient
		user      *user.User
		resLogins []string
		resErr    error
	}{
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != bob.Name()+bob.Surname() {
						t.Errorf("wrong token %s, expected: %s", token, bob.Name()+bob.Surname())
					}
					return bob.ID(), bob.Login(), nil
				},
			},
			repo: &mockRepo{
				getFollowedLogins: func(id int) ([]string, error) {
					if id != bob.ID() {
						t.Errorf("wrong id %d, expected: %d", id, bob.ID())
					}
					return []string{"ben"}, nil
				},
			},
			user:      bob,
			resLogins: []string{"ben"},
			resErr:    nil,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != ben.Name()+ben.Surname() {
						t.Errorf("wrong token %s, expected: %s", token, ben.Name()+ben.Surname())
					}
					return 0, "", common.ErrInvalidToken
				},
			},
			user:      ben,
			resLogins: nil,
			resErr:    common.ErrInvalidToken,
		},
		{
			client: &mockAuthClient{
				validateToken: func(token string) (int, string, error) {
					if token != bill.Name()+bill.Surname() {
						t.Errorf("wrong token %s, expected: %s", token, bill.Name()+bill.Surname())
					}
					return bill.ID(), bill.Login(), nil
				},
			},
			repo: &mockRepo{
				getFollowedLogins: func(id int) ([]string, error) {
					if id != bill.ID() {
						t.Errorf("wrong id %d, expected: %d", id, bob.ID())
					}
					return nil, errors.New("service error")
				},
			},
			user:      bill,
			resLogins: nil,
			resErr:    common.ErrInternal,
		},
	}

	for _, d := range data {
		service := New(d.repo, d.client)
		logins, err := service.GetFollowedLogins(d.user.Name() + d.user.Surname())
		if err != d.resErr {
			t.Errorf("wrong error %s, expected: %s", err, d.resErr)
		}

		if len(logins) != len(d.resLogins) {
			t.Errorf("wrong logins %q, expected: %q", logins, d.resLogins)
			return
		}

		for i, v := range logins {
			if v != d.resLogins[i] {
				t.Errorf("wrong logins %q, expected: %q", logins, d.resLogins)
				return
			}
		}
	}
}
