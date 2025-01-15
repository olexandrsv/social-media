package repository

type MessageModel struct {
	ID         string   `bson:"_id,omitempty"`
	UserID     int      `bson:"userId"`
	Text       string   `bson:"text"`
	ImagesPath []string `bson:"images"`
	FilesPath  []string `bson:"files"`
}

type CreateMessageReq struct {
	UserID      int
	Text        string
	ImagesPaths []string
	FilesPaths  []string
}

type UpdateMessageModel struct {
	Text       string   `bson:"text"`
	ImagesPath []string `bson:"images"`
	FilesPath  []string `bson:"files"`
}

type CreatePostReq struct {
	CreateMessageReq `bson:",inline"`
}

type PostModel struct {
	MessageModel `bson:",inline"`
	CommentsIDs  []string `bson:"comments"`
}

type UpdatePostModel struct {
	UpdateMessageModel
}

type CreatePostCommentReq struct {
	PostID string
	CreateMessageReq
}

type CommentModel struct {
	MessageModel `bson:",inline"`
	CommentsIDs  []string `bson:"comments"`
}

type CreateCommentReq struct {
	CreateMessageReq
}

type CreateCommentCommentReq struct {
	CommentID string
	CreateMessageReq
}

type UpdateCommentModel struct {
	UpdateMessageModel `bson:",inline"`
}
