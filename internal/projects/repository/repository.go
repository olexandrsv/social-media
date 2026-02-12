package repository

import (
	"database/sql"
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/projects/domain/project"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository interface {
	CreateProject(string, string, string, int) (int, error)
	UpdateProject(*project.Project) error
	GetProjects() ([]*project.Project, error)
	GetProject(int) (*project.Project, error)
	DeleteProject(int) error
}

type repo struct {
	db *sql.DB
}

func New() Repository {
	user := config.App.Projects.DB.User
	password := config.App.Projects.DB.Password
	host := config.App.Projects.DB.Host
	port := config.App.Projects.DB.Port
	name := config.App.Projects.DB.Name
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error(err)
		panic("Unable to connect to database")
	}
	return &repo{
		db: db,
	}
}

func (r *repo) CreateProject(name, description, stack string, userID int) (int, error) {
	sql := "INSERT INTO projects (title,description,stack,user_id) VALUES ($1,$2,$3,$4) RETURNING id"
	row := r.db.QueryRow(sql, name, description, stack, userID)

	var id int
	if err := row.Scan(&id); err != nil {
		log.Error(errors.WithStack(err))
		return 0, common.ErrInternal
	}
	return id, nil
}

func (r *repo) UpdateProject(project *project.Project) error {
	sql := "UPDATE projects SET title = $1, description = $2, stack = $3 WHERE id = $4"
	_, err := r.db.Exec(sql, project.Name(), project.Description(), project.Stack(), project.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) GetProject(id int) (*project.Project, error) {
	query := "SELECT id,title,description,stack,user_id FROM projects WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var p ProjectModel
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Stack, &p.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrNotFound
		}
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	return project.New(p.ID, p.Name, p.Description, p.Stack, p.UserID), nil
}

func (r *repo) GetProjects() ([]*project.Project, error){
	sql := "SELECT id,title,description,stack,user_id FROM projects"
	rows, err := r.db.Query(sql)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	
	var projects []*project.Project
	for rows.Next() {
		var p ProjectModel
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Stack, &p.UserID); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		project := project.New(p.ID, p.Name, p.Description, p.Stack, p.UserID)
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *repo) DeleteProject(id int) error {
	sql := "DELETE FROM projects WHERE id = $1"
	_, err := r.db.Exec(sql, id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}