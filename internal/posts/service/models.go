package service

import "mime/multipart"

type UpdatePostReq struct {
	Token         string
	ID            string
	Text          string
	Images        []*multipart.FileHeader
	Files         []*multipart.FileHeader
	DeletedImages []string
	DeletedFiles  []string
}

type CreatePostCommentReq struct {
	Token  string
	PostID string
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type UpdateCommentReq struct{
	Token         string
	ID            string
	Text          string
	Images        []*multipart.FileHeader
	Files         []*multipart.FileHeader
	DeletedImages []string
	DeletedFiles  []string
}
