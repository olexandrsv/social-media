package user

type User struct {
	id       int
	login string
}

func New(id int, login string) *User {
	return &User{
		id:    id,
		login: login,
	}
}

func (u *User) ID() int {
	return u.id
}
func (u *User) Login() string {
	return u.login
}

func (u *User) SetLogin(login string) {
	u.login = login
}