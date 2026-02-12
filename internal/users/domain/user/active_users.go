package user

import "sync"

var activeUsers *users

type users struct {
	data map[int]*User
	mu    sync.Mutex
}

func Get(id int) (*User, bool){
	if activeUsers == nil {
		return nil, false
	}
	user, exist := activeUsers.data[id]
	return user, exist
}

func Set(user *User){
	if activeUsers == nil {
		activeUsers = &users{
			data: make(map[int]*User),
		}
	}
	activeUsers.data[user.ID()] = user
}

