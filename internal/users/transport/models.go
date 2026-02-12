package transport

type NotificationType string

const(
	Message NotificationType = "message"
	Post NotificationType = "post"
)

type Notification struct {
	Type NotificationType `json:"type"`
}

type MessageNotification struct {
	ChatID int
	UserID int
	Text   string
	Images []string
	FIles  []string
}

type PostNotification struct {
	UserID int      `json:"user_id"`
	Text   string   `json:"text"`
	Images []string `json:"images"`
	Files  []string `json:"files"`
}

type UpdateRedMessagesReq struct{}

type ReadMessagesNotification struct {
	ChatID            int    `json:"chat_id"`
	LastReadMessageID string `json:"last_read_message_id"`
}

type ReadPostsNotification struct {
	OwnerID        int    `json:"owner_id"`
	LastReadPostID string `json:"last_read_post_id"`
}