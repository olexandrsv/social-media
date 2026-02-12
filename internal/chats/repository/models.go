package repository

type ChatModel struct {
	ID       int
	Name     string
	OwnerID  int
	UsersIDs []int
}

type CreateChatReq struct {
	Name     string
	OwnerID  int
	UsersIDs []int
}

type MessageModel struct {
	ID     string
	ChatID int
}
