package repository

type UserModel struct {
	ID        int
	Login     string
	Name      string
	Surname   string
	Password  string
	Bio       string
	Interests string
}

func NewUserModel(login, name, surname, password string) UserModel {
	return UserModel{
		Login:    login,
		Name:     name,
		Surname:  surname,
		Password: password,
	}
}

type PostModel struct {
	ID     string
	UserID int
}

type MessageModel struct {
	ID     string
	ChatID int
}
