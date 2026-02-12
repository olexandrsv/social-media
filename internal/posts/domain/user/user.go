package user

type User struct {
	id      int
	name    string
	surname string
}

func New(id int, name, surname string) *User {
	return &User{
		id:      id,
		name:    name,
		surname: surname,
	}
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Surname() string {
	return u.surname
}

func (u *User) AddFullName(name, surname string){
	u.name = name
	u.surname = surname
}