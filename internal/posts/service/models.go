package service

import "mime/multipart"

type CreateMessageReq struct{
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type UpdateMessageReq struct{
	ID            string
	Text          string
	Images        []*multipart.FileHeader
	Files         []*multipart.FileHeader
	DeletedImages []string
	DeletedFiles  []string
}

type UpdatePostReq struct {
	Token         string
	UpdateMessageReq
}

type CreatePostCommentReq struct {
	Token  string
	PostID string
	CreateMessageReq
}

type UpdateCommentReq struct{
	Token         string
	UpdateMessageReq
}

type CreateCommentCommentReq struct{
	Token string
	ParentID string
	CreateMessageReq
}