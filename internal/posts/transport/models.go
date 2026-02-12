package transport

import (
	"mime/multipart"
	"social-media/internal/posts/domain/chatmessage"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/models"
)

type CreatePostResp models.PostModel

type UpdatePostResp models.PostModel

type GetPostsResp []models.PostModel

type PostCommentsResp struct {
	Tone     ToneModel      `json:"tone"`
	Comments []CommentModel `json:"comments"`
}

type CommentsReq struct {
	Token    string
	ParentID string
}

type CreateMessageReq struct {
	Text       string
	FilesPath  []string
	ImagesPath []string
}

type UpdateMessageReq struct {
	Token         string
	ID            string
	Text          string
	Files         []*multipart.FileHeader
	Images        []*multipart.FileHeader
	DeletedFiles  []string
	DeletedImages []string
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

type ToneModel struct {
	PositivePercentage float64 `json:"positive_percentage"`
	NegativePercentage float64 `json:"negative_percentage"`
}

type CommentModel struct {
	models.MessageModel
}

func commentToModel(c *comment.Comment) CommentModel {
	return CommentModel{
		MessageModel: models.MessageToModel(c.Message),
	}
}

type CreatePostCommentReq struct {
	Token  string
	PostID string
	Text   string
	Images []*multipart.FileHeader
	Files  []*multipart.FileHeader
}

type DeletePostCommentReq struct {
	Token     string
	ParentID  string
	CommentID string
}

type UpdateCommentReq struct {
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

type DeleteCommentCommentReq struct {
	Token     string
	ParentID  string
	CommentID string
}

type Token struct {
	Token string
}

type ChatMessageModel struct {
	models.MessageModel
	ChatID int `json:"chat_id"`
}

func chatMessageToModel(m *chatmessage.ChatMessage) ChatMessageModel {
	return ChatMessageModel{
		MessageModel: models.MessageToModel(m.Message),
		ChatID:       m.ChatID(),
	}
}

type MissedPostsNumber struct {
	FollowingID int `json:"following_id"`
	Number      int `json:"number"`
}

type MissedMessagesNumber struct {
	ChatID int `json:"chat_id"`
	Number int `json:"number"`
}
