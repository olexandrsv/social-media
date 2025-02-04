package service

import (
	"social-media/internal/chats/domain/chat"
	"social-media/internal/chats/domain/user"
	"social-media/internal/chats/repository"
	"social-media/internal/common"
	"social-media/internal/common/clients"
	"social-media/internal/common/slice"
)

type Service interface {
	UserChats(string) ([]*chat.Chat, error)
	CreateChat(CreateChatReq) (*chat.Chat, error)
	Chat(string, int) (*chat.Chat, error)
	UpdateChat(UpdateChatReq) error
	DeleteChat(string, int) error
}

type service struct {
	auth  clients.AuthClient
	users clients.UsersClient
	repo  repository.Repository
}

func New(repo repository.Repository, auth clients.AuthClient, users clients.UsersClient) Service {
	return &service{
		auth:  auth,
		users: users,
		repo:  repo,
	}
}

func (s *service) UserChats(token string) ([]*chat.Chat, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.repo.UserChats(id)
}

func (s *service) CreateChat(req CreateChatReq) (*chat.Chat, error) {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateChat(repository.CreateChatReq{
		Name:     req.Name,
		OwnerID:  id,
		UsersIDs: req.UsersIDs,
	})
}

func (s *service) Chat(token string, chatID int) (*chat.Chat, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	chat, err := s.repo.Chat(chatID)
	if err != nil {
		return nil, err
	}

	ids := slice.MustConvert(chat.Users(), func(u *user.User) int {
		return u.ID()
	})

	infos, err := s.users.UsersInfo(ids)
	if err != nil {
		return nil, err
	}

	for i, user := range chat.Users() {
		info := infos[i]
		user.AddInfo(info.Login, info.Name, info.Surname)
	}

	return chat, nil
}

func (s *service) UpdateChat(req UpdateChatReq) error {
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return err
	}

	ownerID, err := s.repo.ChatOwner(req.ID)
	if err != nil {
		return err
	}

	if ownerID != id {
		return common.ErrForbidden
	}

	users := slice.MustConvert(req.UsersIDs, func(id int) *user.User {
		return user.New(id)
	})
	c := chat.New(req.ID, req.Name, user.New(id), users)

	return s.repo.UpdateChat(c)
}

func (s *service) DeleteChat(token string, chatID int) error {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	ownerID, err := s.repo.ChatOwner(chatID)
	if err != nil {
		return err
	}

	if ownerID != id {
		return common.ErrForbidden
	}

	return s.repo.DeleteChat(chatID)
}