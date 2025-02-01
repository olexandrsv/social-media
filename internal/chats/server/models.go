package server

type UserChatsResp []BasicUserChat

type BasicUserChat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
