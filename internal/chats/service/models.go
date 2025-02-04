package service

type CreateChatReq struct {
	Token    string
	Name     string
	UsersIDs []int
}
