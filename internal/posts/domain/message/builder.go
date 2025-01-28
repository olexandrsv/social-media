package message

type Builder struct {
	m *Message
}

func NewBuilder(messageID string) *Builder {
	return &Builder{
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
	b.m.userID = userID
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
	return b.m
}
