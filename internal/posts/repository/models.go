package repository

type PostModel struct {
	UserID     int
	Text       string
	CommentsID []string
	ImagesPath []string
	FilesPath  []string
}
