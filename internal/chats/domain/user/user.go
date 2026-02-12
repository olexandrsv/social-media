package user

type User struct{
	id int
	login string
	name string
	surname string
}

func New(id int) *User {
	return &User{
		id: id,
	}
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Login() string{
	return u.login
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Surname() string {
	return u.surname
}

func (u *User) AddInfo(login, name, surname string) {
	u.login = login
	u.name = name
	u.surname = surname
}