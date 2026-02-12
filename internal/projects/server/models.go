package server

type ProjectModel struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stack       string `json:"stack"`
}

type CreateProjectResp struct {
	ID int `json:"id"`
}
