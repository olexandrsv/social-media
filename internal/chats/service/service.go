package service

import (
	"social-media/internal/chats/domain/chat"
	"social-media/internal/chats/repository"
	"social-media/internal/common/clients"
)

type Service interface {
	UserChats(string) ([]*chat.Chat, error)
	CreateChat(CreateChatReq) (*chat.Chat, error)
}

type service struct {
	auth  clients.AuthClient
	users clients.UsersClient
	repo repository.Repository
}

func New(repo repository.Repository, auth  clients.AuthClient, users clients.UsersClient) Service {
	return &service{
		auth: auth,
		users: users,
		repo: repo,
	}
}

func (s *service) UserChats(token string) ([]*chat.Chat, error) {
	id, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.repo.UserChats(id)
}

func (s *service) CreateChat(req CreateChatReq) (*chat.Chat, error){
	id, _, err := s.auth.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateChat(repository.CreateChatReq{
		Name: req.Name,
		OwnerID: id,
		UsersIDs: req.UsersIDs,
	})
}
