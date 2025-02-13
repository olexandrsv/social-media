package service

import (
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/files"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"
	"social-media/internal/posts/repository"

	"github.com/pkg/errors"
)

type chatMessagesService interface {
	ChatMessages(string, int) ([]*chatmessage.ChatMessage, error)
	ChatMessage(string, string) (*chatmessage.ChatMessage, error)
	CreateChatMessage(CreateChatMessageReq) (*chatmessage.ChatMessage, error)
	UpdateChatMessage(UpdateChatMessageReq) error
	DeleteChatMessage(string, string) error
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

	log.Logf("messages: %#v", messages)

	ids := slice.MustConvert(messages, func(m *chatmessage.ChatMessage) int {
		return m.UserID()
	})

	log.Logf("ids: %#v", ids)

	infos, err := s.users.UsersInfo(ids)
	if err != nil {
		return nil, err
	}
	log.Logf("%#v", infos)


	for i, message := range messages {
		message.AddUserFullName(infos[i].Name, infos[i].Surname)
	}
	return messages, nil
}

func (s *postsService) ChatMessage(token string, messageID string) (*chatmessage.ChatMessage, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.repo.ChatMessage(messageID)
}

func (s *postsService) CreateChatMessage(req CreateChatMessageReq) (*chatmessage.ChatMessage, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidToken
	}

	images, err := files.Process(req.Images)
	if err != nil {
		return nil, err
	}
	files, err := files.Process(req.Files)
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

	return m, nil
}

func (s *postsService) UpdateChatMessage(req UpdateChatMessageReq) error {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidToken
	}

	message, err := s.repo.ChatMessage(req.ID)
	if err != nil {
		return err
	}

	if message.UserID() != id {
		return common.ErrForbidden
	}

	err = message.Update(req.Text, req.Images, req.Files, req.DeletedImages, req.DeletedFiles)
	if err != nil {
		return err
	}

	return s.repo.UpdateChatMessage(message)
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

	return s.repo.DeleteChatMessage(messageID)
}
