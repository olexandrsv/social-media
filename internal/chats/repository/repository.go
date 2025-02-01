package repository

import (
	"database/sql"
	"fmt"
	"social-media/internal/chats/domain/chat"
	"social-media/internal/chats/domain/user"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository interface {
	UserChats(int) ([]*chat.Chat, error)
	ChatDetails(int) (*chat.Chat, error)
}

type repo struct {
	db *sql.DB
}

func New() Repository {
	user := config.App.PostgresDB.User
	password := config.App.PostgresDB.Password
	host := config.App.PostgresDB.Host
	port := config.App.PostgresDB.Port
	name := config.App.PostgresDB.Name
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
