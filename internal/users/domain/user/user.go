package user

import (
	"social-media/internal/common/app/log"
	"social-media/internal/common/connection"
)

type Option func(user *User)

type User struct {
	id        int
	login     string
	name      string
	surname   string
	password  string
	bio       string
	interests string
	conn      connection.MessageConnection
}

func New(id int, login string, opts ...Option) *User {
	user := &User{
		id:    id,
		login: login,
	}
	for _, opt := range opts {
		opt(user)
	}
	return user
}

func WithName(name string) Option {
	return func(user *User) { user.name = name }
}

func WithSurname(surname string) Option {
	return func(user *User) { user.surname = surname }
}

func WithBio(bio string) Option {
	return func(user *User) { user.bio = bio }
}

func WithInterests(interests string) Option {
	return func(user *User) { user.interests = interests }
}

func WithPassword(password string) Option {
	return func(user *User) { user.password = password }
}

func (u *User) SetConn(conn connection.MessageConnection) {
	u.conn = conn
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Surname() string {
	return u.surname
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Bio() string {
	return u.bio
}

func (u *User) Interests() string {
	return u.interests
}

func (u *User) SetConnection(conn connection.MessageConnection) {
	u.conn = conn
}

func (u *User) Send(data string) error {
	log.Infof("notification: %v", data)
	return u.conn.Send([]byte(data))
}
