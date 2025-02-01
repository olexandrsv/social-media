package repository

type ChatModel struct {
	ID       int
	Name     string
	OwnerID  int
	UsersIDs []int
}
