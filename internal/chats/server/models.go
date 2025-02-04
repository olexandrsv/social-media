package server

import (
	"social-media/internal/chats/domain/chat"
	"social-media/internal/chats/domain/user"
)

type UserChatsResp []BasicUserChat

type BasicUserChat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserModel struct {
	ID      int    `json:"id,omitempty"`
	Login   string `json:"login,omitempty"`
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

func userToModel(u *user.User) UserModel {
	return UserModel{
		ID:      u.ID(),
		Login:   u.Login(),
		Name:    u.Name(),
		Surname: u.Surname(),
	}
}

type ChatModel struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Owner UserModel   `json:"owner"`
	Users []UserModel `json:"users"`
}

func chatToModel(c *chat.Chat) ChatModel{
	users := make([]UserModel, 0, len(c.Users()))
	for _, user := range c.Users() {
		users = append(users, userToModel(user))
	}

	return ChatModel{
		ID: c.ID(),
		Name: c.Name(),
		Owner: userToModel(c.Owner()),
		Users: users,
	}
}

type CreateChatResp ChatModel
