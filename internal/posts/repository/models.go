package repository

type PostModel struct {
	ID         string   `bson:"_id,omitempty"`
	UserID     int      `bson:"userId"`
	Text       string   `bson:"text"`
	CommentsID []string `bson:"comments"`
	ImagesPath []string `bson:"images"`
	FilesPath  []string `bson:"files"`
}
