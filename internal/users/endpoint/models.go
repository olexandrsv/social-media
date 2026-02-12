package endpoint

type CreateUserReq struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type AuthResp struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserReq struct {
	ID    int
	Token string
}

type GetUserResp struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Bio       string `json:"bio"`
	Interests string `json:"interests"`
}

type UpdateUserReq struct {
	Token     string
	Name      string
	Surname   string
	Bio       string
	Interests string
}

type UpdateUserResp struct {
	Error string `json:"error"`
}

type GetUsersByInfoReq struct {
	Info string
}

type LoginsResp struct {
	Logins []string `json:"logins"`
}

type FollowUserReq struct {
	Token string
	ID    int
}

type TokenReq struct {
	Token string
}

type UsersResp []UserModel

type UserModel struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type UsersInfoReq struct {
	IDs []int
}

type UsersInfoResp struct {
	Info []UserInfo
}

type UserInfo struct {
	Login   string
	Name    string
	Surname string
}
