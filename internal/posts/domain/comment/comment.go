package comment

type Option func(c *Comment)

type Comment struct {
	id          string
	userID      int
	text        string
	imagesPaths []string
	filesPaths  []string
	commentsIDs []string
}

func New(id string, userID int, opts ...Option) *Comment{
	c := &Comment{
		id:     id,
		userID: userID,
	}
	for _, opt := range opts{
		opt(c)
	}
	return c
}

func WithText(text string) Option {
	return func(c *Comment) { c.text = text }
}

func WithImagesPaths(imagesPaths []string) Option {
	return func(c *Comment) { c.imagesPaths = imagesPaths }
}

func WithFilesPaths(filesPaths []string) Option {
	return func(c *Comment) { c.filesPaths = filesPaths }
}

func WithCommentsIDs(commentIDs []string) Option {
	return func(c *Comment) { c.commentsIDs = commentIDs }
}

func (c *Comment) ID() string {
	return c.id
}

func (c *Comment) UserID() int {
	return c.userID
}

func (c *Comment) Text() string {
	return c.text
}

func (c *Comment) ImagesPaths() []string {
	return c.imagesPaths
}

func (c *Comment) FilesPaths() []string {
	return c.filesPaths
}

func (c *Comment) CommentsIDs() []string {
	return c.commentsIDs
}