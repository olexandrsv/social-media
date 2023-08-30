package endpoint

type CreateUserReq struct {
	Login      string `json:"login"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Password   string `json:"password"`
}

type AuthResp struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserReq struct {
	Login string `json:"login"`
}

type GetUserResp struct {
	Login      string `json:"login"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Bio        string `json:"bio"`
	Interests  string `json:"interests"`
}

type UpdateUserReq struct {
	Token      string `json:"token"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Bio        string `json:"bio"`
	Interests  string `json:"interests"`
}

type UpdateUserResp struct {
	Error string `json:"error"`
}

type GetUserByInfoReq struct {
	Info string
}

type LoginsResp struct {
	Logins []string `json:"logins"`
}

type FollowUserReq struct {
	Token string
	Login string
}

type TokenReq struct{
	Token string
}

