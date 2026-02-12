package service

import (
	"social-media/internal/common"
	"social-media/internal/common/clients"
	"social-media/internal/common/slice"
	"social-media/internal/projects/domain/project"
	"social-media/internal/projects/repository"
)

type Service interface {
	CreateProject(string, string, string, string) (int, error)
	UpdateProject(string, int, string, string, string) error
	GetProjects(string) ([]*project.Project, error)
	DeleteProject(string, int) error
}

type service struct {
	auth clients.AuthClient
	users clients.UsersClient
	repo repository.Repository
}

func New(repo repository.Repository, auth clients.AuthClient, users clients.UsersClient) Service {
	return &service{
		repo: repo,
		users: users,
		auth: auth,
	}
}

func (s *service) CreateProject(token, name, description, stack string) (int, error) {
	userID, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateProject(name, description, stack, userID)
}

func (s *service) UpdateProject(token string, projectID int, name, description, stack string) error {
	userID, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	p, err := s.repo.GetProject(projectID)
	if err != nil {
		return err
	}

	if p.UserID() != userID {
		return common.ErrForbidden
	}

	return s.repo.UpdateProject(project.New(projectID, name, description, stack, userID))
}

func (s *service) GetProjects(token string) ([]*project.Project, error) {
	_, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	projects, err := s.repo.GetProjects()
	if err != nil {
		return nil, err
	}

	usersIDs := slice.MustConvert(projects, func(p *project.Project) int {
		return p.UserID()
	})
	fullNames, err := s.users.UsersInfo(usersIDs)
	if err != nil {
		return nil, err
	}

	for i, fullName := range fullNames {
		projects[i].SetUserLogin(fullName.Login)
	}

	return projects, nil
}

func (s *service) DeleteProject(token string, projectID int) error {
	userID, _, err := s.auth.ValidateToken(token)
	if err != nil {
		return err
	}

	p, err := s.repo.GetProject(projectID)
	if err != nil {
		return err
	}

	if p.UserID() != userID {
		return common.ErrForbidden
	}

	return s.repo.DeleteProject(projectID)
}
