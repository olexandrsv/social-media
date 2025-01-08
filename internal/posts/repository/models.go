package repository

type PostModel struct {
	ID          string   `bson:"_id,omitempty"`
	UserID      int      `bson:"userId"`
	Text        string   `bson:"text"`
	CommentsIDs []string `bson:"comments"`
	ImagesPath  []string `bson:"images"`
	FilesPath   []string `bson:"files"`
}

type UpdatePostModel struct {
	Text       string   `bson:"text"`
	ImagesPath []string `bson:"images"`
	FilesPath  []string `bson:"files"`
}

type CommentModel struct {
	ID          string   `bson:"_id,omitempty"`
	UserID      int      `bson:"userId"`
	Text        string   `bson:"text"`
	CommentsIDs []string `bson:"comments"`
	ImagesPath  []string `bson:"images"`
	FilesPath   []string `bson:"files"`
}

type CreateCommentReq struct {
	UserID     int
	PostID     string
	Text       string
	ImagesPath []string
	FilesPath  []string
}

type UpdateCommentModel struct {
	Text       string   `bson:"text"`
	ImagesPath []string `bson:"images"`
	FilesPath  []string `bson:"files"`
}
