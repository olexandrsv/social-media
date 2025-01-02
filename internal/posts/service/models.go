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
