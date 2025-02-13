package service

import "mime/multipart"

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

type UpdateChatMessageReq struct{
	UpdateMessageReq
}
