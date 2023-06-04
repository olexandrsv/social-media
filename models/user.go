package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

var ActiveUsers = userMap{
	data: make(map[int]*User),
}

type userMap struct {
	mux  sync.RWMutex
	data map[int]*User
}

func (users userMap) Get(key int) (*User, bool) {
	users.mux.RLock()
	defer users.mux.RUnlock()
	user, ok := users.data[key]
	return user, ok
}

func (users userMap) Set(key int, user *User) {
	users.mux.Lock()
	defer users.mux.Unlock()
	users.data[key] = user
}

type User struct {
	Id    int
	Login string
	Conn  *websocket.Conn
}
