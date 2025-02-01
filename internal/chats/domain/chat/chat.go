package chat

import "social-media/internal/chats/domain/user"

type Chat struct {
	id    int
	name  string
	owner *user.User
	users []*user.User
}

func New(id int, name string, owner *user.User, users []*user.User) *Chat {
	return &Chat{
		id: id,
		name: name,
		owner: owner,
		users: users,
	}
}

func (c *Chat) ID() int {
	return c.id
}

func (c *Chat) Name() string {
	return c.name
}

func (c *Chat) Owner() *user.User {
	return c.owner
}

func (c *Chat) Users() []*user.User {
	return c.users
}
