package endpoint

type CreatePostReq struct {
	Token      string
	Text       string
	FilesPath  []string
	ImagesPath []string
}

type CreatePostResp struct {
	ID          string   `json:"_id"`
	Text        string   `json:"text"`
	FilesPaths  []string `json:"files"`
	ImagesPaths []string `json:"images"`
}
