package post

type Post struct {
	ID     string
	UserID int
}

func NewPost(id string, userID int) *Post {
	return &Post{
		ID:     id,
		UserID: userID,
	}
}
