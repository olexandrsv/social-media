package comment

import (
	"social-media/internal/posts/domain/message"
)

type Comment struct {
	*message.Message
	commentsIDs []string
}

func New(id string, userID int, text string, imagesPaths, filesPaths, commentsIDs []string) *Comment {
	m := message.New(id, userID, text, imagesPaths, filesPaths)
	return &Comment{
		Message:     m,
		commentsIDs: commentsIDs,
	}
}

func (c *Comment) CommentsIDs() []string {
	return c.commentsIDs
}