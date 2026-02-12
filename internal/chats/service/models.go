package service

type CreateChatReq struct {
	Token    string
	Name     string
	UsersIDs []int
}

type UpdateChatReq struct {
	Token    string
	ID       int
	Name     string
	UsersIDs []int
}
