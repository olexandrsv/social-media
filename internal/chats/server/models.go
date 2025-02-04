package server

import "social-media/internal/chats/domain/user"

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
	Name  string      `json:"name"`
	Owner UserModel   `json:"owner"`
	Users []UserModel `json:"users"`
}

type CreateChatResp ChatModel
