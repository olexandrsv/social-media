package models

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	"social-media/hash"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	ID         int             `json:"id,omitempty"`
	Login      string          `json:"login,omitempty"`
	FirstName  string          `json:"firstName,omitempty"`
	SecondName string          `json:"secondName,omitempty"`
	Password   string          `json:"-"`
	Bio        string          `json:"bio,omitempty"`
	Interests  string          `json:"interests"`
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

func ParseToken(token string) (*User, error) {
	id, login, err := auth.TokenCredentials(token)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:    id,
		Login: login,
	}, nil
}

// Save stores user in PostgreSQL database
func (user *User) Save() error {
	var id int
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "insert into users (login, first_name, second_name, password, bio, interests) values ($1, $2, $3, $4, $5, $6) returning id", user.Login, user.FirstName, user.SecondName, user.Password, "", "").Scan(&id)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// Update updates user's data in PostgreSQL database and in ActiveUsers map
func (user *User) Update() error{
	conn := database.PostgreConn
	_, err := conn.Exec(context.Background(), "update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5", user.FirstName, user.SecondName, user.Bio, user.Interests, user.ID)
	if err != nil{
		return err
	}
	if prevUser, ok := ActiveUsers.Get(user.ID); ok{
		user.Conn = prevUser.Conn
		ActiveUsers.Set(user.ID, user)
	}
	return err
}

// Loads user from PostgreSQL database
func Load(login string) (*User, error) {
	var user User
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select first_name, second_name, bio, interests from users where login=$1", login).Scan(&user.FirstName, &user.SecondName, &user.Bio, &user.Interests)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Register adds current user to active users
func (user *User) Register() {
	ActiveUsers.Set(user.ID, user)
}

func (user *User) GenerateJWT() (string, error) {
	return auth.GenerateJWT(user.ID, user.Login)
}

func (user *User) Exist() (bool, error) {
	var encodedPassw string

	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select id, password from users where login=$1", user.Login).Scan(&user.ID, &encodedPassw)
	if err != nil {
		return false, err
	}

	if !hash.CheckPassword(user.Password, encodedPassw) {
		return false, nil
	}
	user.Password = encodedPassw

	return true, nil
}

func (user *User) Subscribe(followedID int) error {
	conn := database.PostgreConn
	_, err := conn.Exec(context.Background(), "insert into followers (user_id, follower_id) values ($1, $2)", followedID, user.ID)
	return err
}

// Followed returns users who are followed by current user
func (user *User) Followed() ([]*User, error){
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select id, login from users join followers on users.id = followers.user_id and followers.follower_id=$1", user.ID)
	if err != nil {
		return nil, err
	}

	var users []*User
	for rows.Next() {
		var followed User
		err = rows.Scan(&followed.ID, &followed.Login)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &followed)
	}
	return users, nil
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

func GetByInfo(info string) (string, error){
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select login from users where interests like $1 or bio like $1", "%"+info+"%")
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for rows.Next() {
		var login string
		err = rows.Scan(&login)
		if err != nil {
			continue
		}
		sb.WriteString(",")
		sb.WriteString(login)
	}
	res := sb.String()
	if len(res) > 0 {
		res = res[1:]
	}
	return res, nil
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
