package project

import "social-media/internal/projects/domain/user"

type Project struct {
	id          int
	name        string
	description string
	stack       string
	user        *user.User
}

func New(id int, name, description, stack string, userID int) *Project {
	return &Project{
		id:          id,
		name:        name,
		description: description,
		stack:       stack,
		user: user.New(userID, ""),
	}
}

func (p *Project) ID() int {
	return p.id
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) Description() string {
	return p.description
}

func (p *Project) Stack() string {
	return p.stack
}

func (p *Project) UserID() int {
	return p.user.ID()
}

func (p *Project) UserLogin() string {
	return p.user.Login()
}

func (p *Project) SetUserLogin(login string) {
	p.user.SetLogin(login)
}