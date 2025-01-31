package user

type Builder struct{
	user *User
}

func NewBuilder() *Builder {
	return &Builder{
		user: &User{},
	}
}

func (b *Builder) WithID(id int) *Builder {
	b.user.id = id
	return b
}

func (b *Builder) WithName(name string) *Builder{
	b.user.name = name
	return b
}

func (b *Builder) WithSurname(surname string) *Builder{
	b.user.surname = surname
	return b
}

func (b *Builder) Create() *User{
	return b.user
}
