package endpoint

import "mime/multipart"

type Token struct {
	Token string
}

type GetPostsRequest struct {
	Token  string
	UserID int
}

type CreatePostReq struct {
	Token      string
	Text       string
	FilesPath  []string
	ImagesPath []string
}

type PostModel struct {
	ID          string   `json:"id"`
	UserID      int      `json:"user_id"`
	Text        string   `json:"text"`
	FilesPaths  []string `json:"files"`
	ImagesPaths []string `json:"images"`
}

type UpdatePostReq struct {
	Token         string
	ID            string
	Text          string
	Files         []*multipart.FileHeader
	Images        []*multipart.FileHeader
	DeletedFiles  []string
	DeletedImages []string
}

type DeletePostReq struct {
	Token  string
	PostID string
}

type GetPostCommentsReq struct {
	Token  string
	PostID string
}

type CommentModel struct {
	ID         string   `json:"_id,omitempty"`
	UserID     int      `json:"userId"`
	Text       string   `json:"text"`
	ImagesPath []string `json:"images"`
	FilesPath  []string `json:"files"`
}

type CreatePostCommentReq struct {
	Token  string
	PostID string
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type UpdateCommentReq struct {
	Token         string
	CommentID     string
	Text          string
	Images        []*multipart.FileHeader
	Files         []*multipart.FileHeader
	DeletedFiles  []string
	DeletedImages []string
}
