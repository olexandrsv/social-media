package endpoint

import (
	"mime/multipart"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/domain/post"
)

type MessageModel struct {
	ID          string   `json:"id"`
	UserID      int      `json:"user_id"`
	Text        string   `json:"text"`
	FilesPaths  []string `json:"files"`
	ImagesPaths []string `json:"images"`
}

type CreateMessageReq struct {
	Text       string
	FilesPath  []string
	ImagesPath []string
}

type UpdateMessageReq struct {
	ID            string
	Text          string
	Files         []*multipart.FileHeader
	Images        []*multipart.FileHeader
	DeletedFiles  []string
	DeletedImages []string
}

type PostModel struct {
	MessageModel
}

func postToModel(p *post.Post) PostModel {
	return PostModel{
		MessageModel: MessageModel{
			ID:          p.ID(),
			UserID:      p.UserID(),
			Text:        p.Text(),
			ImagesPaths: p.ImagesPaths(),
			FilesPaths:  p.FilesPaths(),
		},
	}
}

type GetPostsRequest struct {
	Token  string
	UserID int
}

type CreatePostReq struct {
	Token string
	CreateMessageReq
}

type UpdatePostReq struct {
	Token string
	UpdateMessageReq
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
	MessageModel
}

func commentToModel(c *comment.Comment) CommentModel {
	return CommentModel{
		MessageModel: MessageModel{
			ID:          c.ID(),
			UserID:      c.UserID(),
			Text:        c.Text(),
			ImagesPaths: c.ImagesPaths(),
			FilesPaths:  c.FilesPaths(),
		},
	}
}

type CreatePostCommentReq struct {
	Token  string
	PostID string
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type UpdateCommentReq struct {
	Token string
	UpdateMessageReq
}

type GetCommentCommentsReq struct {
	Token     string
	CommentID string
}

type CreateCommentCommentReq struct {
	Token    string
	ParentID string
	Text     string
	Images   []*multipart.FileHeader
	Files    []*multipart.FileHeader
}

type Token struct {
	Token string
}
