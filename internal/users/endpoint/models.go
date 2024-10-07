package endpoint

type CreateUserReq struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type AuthResp struct {
	Token string `json:"token"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserReq struct {
	Login string `json:"login"`
}

type GetUserResp struct {
	Login     string `json:"login"`
	Name      string `json:"first_name"`
	Surname   string `json:"second_name"`
	Bio       string `json:"bio"`
	Interests string `json:"interests"`
}

type UpdateUserReq struct {
	Token     string `json:"token"`
	Name      string `json:"first_name"`
	Surname   string `json:"second_name"`
	Bio       string `json:"bio"`
	Interests string `json:"interests"`
}

type UpdateUserResp struct {
	Error string `json:"error"`
}

type GetLoginsByInfoReq struct {
	Info string
}

type LoginsResp struct {
	Logins []string `json:"logins"`
}

type FollowUserReq struct {
	Token string
	Login string
}

type TokenReq struct {
	Token string
}
