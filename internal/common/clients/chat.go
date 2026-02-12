package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"social-media/api/pb/chats"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type UserModel struct {
	ID      int    `json:"id,omitempty"`
	Login   string `json:"login,omitempty"`
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}

type ChatModel struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Owner UserModel   `json:"owner"`
	Users []UserModel `json:"users"`
}

type ChatsClient interface {
	Chat(string, int) (*ChatModel, error)
	LastReadMessages(userID int) ([]Message, error)
	UpdateMissedMessages(token string, chatID int, lastReadMessage string) error
}

type chatsClient struct {
	chats.ChatsClient
}

func NewChatsClient() ChatsClient {
	host := config.App.Chats.Service.Host
	port := config.App.Chats.Service.GrpcPort
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(err)
	}
	return &chatsClient{
		chats.NewChatsClient(conn),
	}
}

type Message struct {
	ID     string
	ChatID int
}

func (c *chatsClient) LastReadMessages(userID int) ([]Message, error) {
	resp, err := c.ChatsClient.LastReadMessages(context.Background(), &chats.LastReadMessagesReq{
		UserID: int64(userID),
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if resp.Err != nil {
		err = common.NewError(int(resp.Err.Code), resp.Err.Message)
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	convertedMessages := slice.MustConvert(resp.Messages, func(m *chats.Message) Message {
		return Message{
			ID:     m.LastMessageID,
			ChatID: int(m.ChatID),
		}
	})
	return convertedMessages, nil
}

func (c *chatsClient) UpdateMissedMessages(token string, chatID int, lastReadMessage string) error {
	body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
	writer.WriteField("last_read_message", lastReadMessage)
	host := config.App.Chats.Service.Host
	port := config.App.Chats.Service.HttpPort
	url := fmt.Sprintf("http://%s:%s/users/chats/%d/read", host, port, chatID)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return err
}

func (c *chatsClient) Chat(token string, chatID int) (*ChatModel, error) {
	body := &bytes.Buffer{}
	host := config.App.Chats.Service.Host
	port := config.App.Chats.Service.HttpPort
	url := fmt.Sprintf("http://%s:%s/chats/%d", host, port, chatID)
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var model ChatModel
	log.Info(string(data))
	if err := json.Unmarshal(data, &model); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return &model, nil
}
