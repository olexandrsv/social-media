package transport

import (
	"mime/multipart"
	"social-media/internal/posts/domain/chatmessage"
	"social-media/internal/posts/domain/comment"
	"social-media/internal/posts/domain/message"
	"social-media/internal/posts/domain/post"
)

type CreatePostResp PostModel

type UpdatePostResp PostModel

type GetPostsResp []PostModel

type PostCommentsResp []CommentModel

type CommentsReq struct {
	Token    string
	ParentID string
}

type User struct {
	ID      int    `json:"user_id"`
	Name    string `json:"user_name"`
	Surname string `json:"user_surname"`
}

type MessageModel struct {
	ID string `json:"id"`
	User
	Text        string   `json:"text"`
	FilesPaths  []string `json:"files"`
	ImagesPaths []string `json:"images"`
}

func messageToModel(m *message.Message) MessageModel {
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

type PostModel struct {
	MessageModel
}

func postToModel(p *post.Post) PostModel {
	return PostModel{
		MessageModel: messageToModel(p.Message),
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
		MessageModel: messageToModel(c.Message),
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
	MessageModel
	ChatID int `json:"chat_id"`
}

func chatMessageToModel(m *chatmessage.ChatMessage) ChatMessageModel {
	return ChatMessageModel{
		MessageModel: messageToModel(m.Message),
		ChatID: m.ChatID(),
	}
}