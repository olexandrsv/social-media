package post

type Option func(post *Post)

type Post struct {
	id          string
	userID      int
	text        string
	imagesPaths []string
	filesPaths  []string
	commentsIDs []string
}

func New(id string, userID int, opts ...Option) *Post {
	post := &Post{
		id:     id,
		userID: userID,
	}
	for _, opt := range opts {
		opt(post)
	}
	return post
}

func WithText(text string) Option {
	return func(p *Post) { p.text = text }
}

func WithImagesPaths(imagesPaths []string) Option {
	return func(p *Post) { p.imagesPaths = imagesPaths }
}

func WithFilesPaths(filesPaths []string) Option {
	return func(p *Post) { p.filesPaths = filesPaths }
}

func WithCommentsIDs(commentsIDs []string) Option{
	return func(p *Post) { p.commentsIDs = commentsIDs }
}

func (p *Post) ID() string {
	return p.id
}

func (p *Post) UserID() int {
	return p.userID
}

func (p *Post) Text() string {
	return p.text
}

func (p *Post) ImagesPaths() []string {
	return p.imagesPaths
}

func (p *Post) FilesPaths() []string {
	return p.filesPaths
}

func (p *Post) CommentsIDs() []string{
	return p.commentsIDs
}
