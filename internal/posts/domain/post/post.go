package post

import "social-media/internal/posts/domain/message"

type Option func(post *Post)

type Post struct {
	*message.Message
	commentsIDs []string
}

func New(id string, userID int, text string, imagesPaths, filesPaths, commentsIDs []string) *Post {
	m := message.New(id, userID, text, imagesPaths, filesPaths)
	return &Post{
		Message: m,
		commentsIDs: commentsIDs,
	}
}

func (p *Post) CommentsIDs() []string {
	return p.commentsIDs
}
