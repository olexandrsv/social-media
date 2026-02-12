package service

import (
	"mime/multipart"
	"social-media/internal/posts/models"
)

type CreateMessageReq struct {
	Token  string
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type UpdateMessageReq struct {
	Token         string
	ID            string
	Text          string
	Images        []*multipart.FileHeader
	Files         []*multipart.FileHeader
	DeletedImages []string
	DeletedFiles  []string
}

type CreatePostReq struct {
	CreateMessageReq
}

type UpdatePostReq struct {
	UpdateMessageReq
}

type CreatePostCommentReq struct {
	PostID string
	CreateMessageReq
}

type UpdateCommentReq struct {
	UpdateMessageReq
}

type CreateCommentCommentReq struct {
	ParentID string
	CreateMessageReq
}

type CreateChatMessageReq struct {
	ChatID int
	CreateMessageReq
}

type UpdateChatMessageReq struct {
	ChatID int
	UpdateMessageReq
}

type PostCreatedNotification struct {
	Type string           `json:"type"`
	Data models.PostModel `json:"data"`
}

type PostUpdatedNotification struct {
	Type string           `json:"type"`
	Data models.PostModel `json:"data"`
}

type PostDeletedNotification struct {
	Type string           `json:"type"`
	Data models.PostModel `json:"data"`
}

type MessageCreatedNotification struct {
	Type string                  `json:"type"`
	Data models.ChatMessageModel `json:"data"`
}

type MessageUpdatedNotification struct {
	Type string                  `json:"type"`
	Data models.ChatMessageModel `json:"data"`
}

type MessageDeletedNotification struct {
	Type string                  `json:"type"`
	Data models.ChatMessageModel `json:"data"`
}

func newPostCreatedNotification(data models.PostModel) PostCreatedNotification {
	return PostCreatedNotification{
		Type: "postCreated",
		Data: data,
	}
}

func newPostUpdatedNotification(data models.PostModel) PostUpdatedNotification {
	return PostUpdatedNotification{
		Type: "postUpdated",
		Data: data,
	}
}

func newPostDeletedNotification(data models.PostModel) PostDeletedNotification {
	return PostDeletedNotification{
		Type: "postDeleted",
		Data: data,
	}
}

func newMessageCreatedNotification(data models.ChatMessageModel) MessageCreatedNotification {
	return MessageCreatedNotification{
		Type: "messageCreated",
		Data: data,
	}
}

func newMessageUpdatedNotification(data models.ChatMessageModel) MessageUpdatedNotification {
	return MessageUpdatedNotification{
		Type: "messageUpdated",
		Data: data,
	}
}

func newMessageDeletedNotification(data models.ChatMessageModel) MessageDeletedNotification {
	return MessageDeletedNotification{
		Type: "messageDeleted",
		Data: data,
	}
}
