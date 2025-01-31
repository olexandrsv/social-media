package message

import "social-media/internal/posts/domain/user"

type Builder struct {
	userBuilder *user.Builder
	m           *Message
}

func NewBuilder(messageID string) *Builder {
	return &Builder{
		userBuilder: user.NewBuilder(),
		m: &Message{
			id: messageID,
		},
	}
}

func (b *Builder) WithText(text string) *Builder {
	b.m.text = text
	return b
}

func (b *Builder) WithUserID(userID int) *Builder {
	b.userBuilder.WithID(userID)
	return b
}

func (b *Builder) WithUserName(name string) *Builder {
	b.userBuilder.WithName(name)
	return b
}

func (b *Builder) WithUserSurname(surname string) *Builder{
	b.userBuilder.WithSurname(surname)
	return b
}

func (b *Builder) WithImagesPaths(imagesPaths []string) *Builder {
	b.m.imagesPaths = imagesPaths
	return b
}

func (b *Builder) WithFilesPaths(filesPaths []string) *Builder {
	b.m.filesPaths = filesPaths
	return b
}

func (b *Builder) Create() *Message {
	b.m.user = b.userBuilder.Create()
	return b.m
}
