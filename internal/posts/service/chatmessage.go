package service

import (
	"encoding/json"
	"os"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/clients"
	"social-media/internal/common/files"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"
	"social-media/internal/posts/models"
	"social-media/internal/posts/repository"

	"github.com/pkg/errors"
)

type chatMessagesService interface {
	ChatMessages(string, int) ([]*chatmessage.ChatMessage, error)
	CreateChatMessage(CreateChatMessageReq) (*chatmessage.ChatMessage, error)
	UpdateChatMessage(UpdateChatMessageReq) (*chatmessage.ChatMessage, error)
	DeleteChatMessage(string, string) error
	MissedMessagesNumber(string) ([]*chatmessage.MissedNumber, error)
	GetChatMessageFile(string, string) (*os.File, error)
}

func (s *postsService) ChatMessages(token string, chatID int) ([]*chatmessage.ChatMessage, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	messages, err := s.repo.ChatMessages(chatID)
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		if err := s.addSignedUrl(message); err != nil {
			return nil, err
		}
	}

	if err := addUsersFullNames(s, messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *postsService) addSignedUrl(message *chatmessage.ChatMessage) error {
	tokens := make([]string, 0, len(message.ImagesPaths()))
	for _, filePath := range message.ImagesPaths() {
		token, err := s.auth.GenerateSignedUrl(filePath)
		if err != nil {
			return err
		}
		tokens = append(tokens, token)
	}
	message.SetImagesPaths(tokens)
	return nil
}

func (s *postsService) CreateChatMessage(req CreateChatMessageReq) (*chatmessage.ChatMessage, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}

	images, err := files.Process("messages", req.Images)
	if err != nil {
		return nil, err
	}
	files, err := files.Process("messages", req.Files)
	if err != nil {
		return nil, err
	}

	m, err := s.repo.CreateChatMessage(repository.CreateChatMessageReq{
		CreateMessageReq: repository.CreateMessageReq{
			UserID:      id,
			Text:        req.Text,
			ImagesPaths: images,
			FilesPaths:  files,
		},
		ChatID: req.ChatID,
	})
	if err != nil {
		return nil, err
	}

	chat, err := s.chats.Chat(req.Token, req.ChatID)
	if err != nil {
		return nil, err
	}

	infos, err := s.users.UsersInfo([]int{id})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if len(infos) < 1 {
		log.Error(errors.New("wrong users info length"))
		return nil, common.ErrInternal
	}

	if err := s.addSignedUrl(m); err != nil {
		return nil, err
	}
	m.AddUserFullName(infos[0].Name, infos[0].Surname)

	notification := newMessageCreatedNotification(models.ChatMessageToModel(m))
	bytes, err := json.Marshal(notification)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	err = s.users.SendNotification(id, slice.MustConvert(chat.Users, func(u clients.UserModel) int {
		return u.ID
	}), string(bytes))
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	if err := s.chats.UpdateMissedMessages(req.Token, req.ChatID, m.ID()); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *postsService) UpdateChatMessage(req UpdateChatMessageReq) (*chatmessage.ChatMessage, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}

	message, err := s.repo.ChatMessage(req.ID)
	if err != nil {
		return nil, err
	}

	if message.UserID() != id {
		return nil, common.ErrForbidden
	}

	err = message.Update(req.Text, req.Images, req.Files, "messages", req.DeletedImages, req.DeletedFiles)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateChatMessage(message); err != nil {
		return nil, err
	}

	chat, err := s.chats.Chat(req.Token, req.ChatID)
	if err != nil {
		return nil, err
	}

	infos, err := s.users.UsersInfo([]int{id})
	if err != nil {
		return nil, err
	}
	if len(infos) < 1 {
		log.Error(errors.New("wrong users info length"))
		return nil, common.ErrInternal
	}
	message.AddUserFullName(infos[0].Name, infos[0].Surname)

	notification := newMessageUpdatedNotification(models.ChatMessageToModel(message))
	bytes, err := json.Marshal(notification)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	err = s.users.SendNotification(id, slice.MustConvert(chat.Users, func(u clients.UserModel) int {
		return u.ID
	}), string(bytes))
	if err != nil {
		return nil, common.ErrInternal
	}

	return message, nil
}

func (s *postsService) DeleteChatMessage(token, messageID string) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidToken
	}

	message, err := s.repo.ChatMessage(messageID)
	if err != nil {
		return err
	}

	if message.UserID() != id {
		return common.ErrForbidden
	}

	if err := s.repo.DeleteChatMessage(messageID); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidToken
	}

	chat, err := s.chats.Chat(token, message.ChatID())
	if err != nil {
		return err
	}

	notification := newMessageDeletedNotification(models.ChatMessageToModel(message))
	bytes, err := json.Marshal(notification)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return s.users.SendNotification(message.UserID(), slice.MustConvert(chat.Users, func(u clients.UserModel) int {
		return u.ID
	}), string(bytes))
}

func (s *postsService) MissedMessagesNumber(token string) ([]*chatmessage.MissedNumber, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	messages, err := s.chats.LastReadMessages(id)
	if err != nil {
		return nil, err
	}
	convertedMessages := slice.MustConvert(messages, func(m clients.Message) *chatmessage.ChatMessage {
		return chatmessage.New(m.ID, 0, "", nil, nil, m.ChatID)
	})

	return s.repo.MissedMessagesNumber(id, convertedMessages)
}

func (s *postsService) GetChatMessageFile(token, signedUrl string) (*os.File, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	fileID, err := s.auth.ValidateSignedUrl(signedUrl)
	if err != nil {
		return nil, err
	}

	return files.Get("messages", fileID)
}
