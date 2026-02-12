package comment

import "social-media/internal/posts/domain/message"

type Builder struct {
	comment *Comment
	message.Builder
}

func NewBuilder(commentID string) *Builder {
	messageBuilder := message.NewBuilder(commentID)
	return &Builder{
		Builder: *messageBuilder,
	}
}

func (b *Builder) WithCommentsIDs(commentsIDs []string) *Builder {
	b.comment.commentsIDs = commentsIDs
	return b
}

func (b *Builder) Create() *Comment{
	message := b.Builder.Create()
	b.comment.Message = message
	return b.comment
}