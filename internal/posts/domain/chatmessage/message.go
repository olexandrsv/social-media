package chatmessage

import "social-media/internal/posts/domain/message"

type ChatMessage struct {
	*message.Message
	chatID int
}

func New(id string, userID int, text string, imagesPaths, filesPaths []string, chatID int) *ChatMessage {
	m := message.New(id, userID, text, imagesPaths, filesPaths)
	return &ChatMessage{
		Message: m,
		chatID: chatID,
	}
}

func (c *ChatMessage) ChatID() int {
	return c.chatID
}