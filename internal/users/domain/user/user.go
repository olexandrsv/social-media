package user

import (
	"context"
	"social-media/database"
	"social-media/hash"
	"sync"

	"github.com/gorilla/websocket"
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
	Conn      *websocket.Conn
}

func New(login string, name string, surname string, password string) (*User, error) {
	hashPsw, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		login:    login,
		name:     name,
		surname:  surname,
		password: hashPsw,
	}, nil
}

func NewUser(id int, login string, opts ...Option) *User {
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

// Register adds current user to active users
func (u *User) Register() {
	ActiveUsers.Set(u.id, u)
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

func GetIdByLogin(login string) (int, error) {
	var id int
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select id from users where login=$1", login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

var ActiveUsers = userMap{
	data: make(map[int]*User),
}

type userMap struct {
	mux  sync.RWMutex
	data map[int]*User
}

func (users *userMap) Get(id int) (*User, bool) {
	users.mux.RLock()
	defer users.mux.RUnlock()
	user, ok := users.data[id]
	return user, ok
}

func (users *userMap) Set(id int, user *User) {
	users.mux.Lock()
	defer users.mux.Unlock()
	users.data[id] = user
}
