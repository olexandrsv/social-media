package models

import (
	"social-media/internal/posts/domain/chatmessage"
	"social-media/internal/posts/domain/message"
	"social-media/internal/posts/domain/post"
)

type User struct {
	ID      int    `json:"user_id"`
	Name    string `json:"user_name"`
	Surname string `json:"user_surname"`
}

type ChatMessageModel struct {
	ChatID       int `json:"chat_id"`
	MessageModel `json:",inline"`
}

func ChatMessageToModel(m *chatmessage.ChatMessage) ChatMessageModel {
	return ChatMessageModel{
		ChatID:       m.ChatID(),
		MessageModel: MessageToModel(m),
	}
}

type MessageModel struct {
	ID string `json:"id"`
	User
	Text        string   `json:"text"`
	FilesPaths  []string `json:"files"`
	ImagesPaths []string `json:"images"`
}

func MessageToModel(m message.MessageI) MessageModel {
	return MessageModel{
		ID: m.ID(),
		User: User{
			ID:      m.UserID(),
			Name:    m.UserName(),
			Surname: m.UserSurname(),
		},
		Text:        m.Text(),
		ImagesPaths: m.ImagesPaths(),
		FilesPaths:  m.FilesPaths(),
	}
}

type PostModel struct {
	MessageModel
}

func PostToModel(p *post.Post) PostModel {
	return PostModel{
		MessageModel: MessageToModel(p.Message),
	}
}
