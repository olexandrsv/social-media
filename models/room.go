package models

import (
	"sync"
)

var ActiveRoom = roomMap{
	data: make(map[int]*Room),
}

type roomMap struct {
	mux  sync.RWMutex
	data map[int]*Room
}

func (rooms *roomMap) Get(key int) (*Room, bool) {
	rooms.mux.RLock()
	defer rooms.mux.RUnlock()
	room, ok := rooms.data[key]
	return room, ok
}

func (rooms *roomMap) Set(key int, value *Room) {
	rooms.mux.Lock()
	defer rooms.mux.Unlock()
	rooms.data[key] = value
}

type Room struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Users []int  `json:"users,omitempty"`
}
