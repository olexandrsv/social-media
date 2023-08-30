package users

import (
	"context"
	"social-media/database"
	"social-media/hash"
	"sync"

	"github.com/gorilla/websocket"
)

type Option func(user *User)

type User struct {
	ID         int             `json:"id,omitempty"`
	Login      string          `json:"login,omitempty"`
	FirstName  string          `json:"firstName,omitempty"`
	SecondName string          `json:"secondName,omitempty"`
	Password   string          `json:"-"`
	Bio        string          `json:"bio,omitempty"`
	Interests  string          `json:"interests,omitempty"`
	Conn       *websocket.Conn `json:"-"`
}

func New(login string, firstName string, secondName string, password string) (*User, error) {
	hashPsw, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Login:      login,
		FirstName:  firstName,
		SecondName: secondName,
		Password:   hashPsw,
	}, nil
}

func NewUser(login string, opts ...Option) *User{
	user := &User{
		Login: login,
	}
	for _, opt := range opts{
		opt(user)
	}
	return user
}

func WithID(id int) Option{
	return func(user *User) {user.ID = id}
}

func WithFirstName(firstName string) Option{
	return func(user *User) {user.FirstName = firstName}
}

func WithSecondName(secondName string) Option{
	return func(user *User) {user.SecondName = secondName}
}

func WithBio(bio string) Option{
	return func(user *User) {user.Bio = bio}
}

func WithInterests(interests string) Option{
	return func(user *User) {user.Interests = interests}
}

func WithPassword(password string) Option{
	return func(user *User) {user.Password = password}
}

// Register adds current user to active users
func (user *User) Register() {
	ActiveUsers.Set(user.ID, user)
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
