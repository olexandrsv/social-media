package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"social-media/internal/chats/domain/chat"
	"social-media/internal/chats/domain/user"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository interface {
	UserChats(int) ([]*chat.Chat, error)
	ChatDetails(int) (*chat.Chat, error)
	CreateChat(CreateChatReq) (*chat.Chat, error)
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

func (r *repo) UserChats(userID int) ([]*chat.Chat, error) {
	query := `SELECT chat_id, name FROM chats INNER JOIN 
		(SELECT chat_id, user_id FROM chats_users WHERE user_id = $1) AS user_chats
		ON chats.id = user_chats.chat_id;`
	rows, err := r.db.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var chats []*chat.Chat
	for rows.Next() {
		var model ChatModel
		if err := rows.Scan(&model.ID, &model.Name); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		c := chat.New(model.ID, model.Name, nil, nil)
		chats = append(chats, c)
	}

	return chats, nil
}

func (r *repo) ChatDetails(id int) (*chat.Chat, error) {
	owner, err := r.getChatOwner(id)
	if err != nil {
		return nil, err
	}
	users, err := r.getChatUsers(id)
	if err != nil {
		return nil, err
	}

	chat := chat.New(id, "", owner, users)
	return chat, nil
}

func (r *repo) getChatOwner(chatID int) (*user.User, error) {
	query := `SELECT owner FROM chats WHERE id=$1`
	row := r.db.QueryRow(query, chatID)

	var id int
	if err := row.Scan(&id); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	return user.New(id), nil
}

func (r *repo) getChatUsers(chatID int) ([]*user.User, error) {
	query := `SELECT user_id FROM chats_users WHERE chat_id=$1`
	rows, err := r.db.Query(query, chatID)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var users []*user.User
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		user := user.New(id)
		users = append(users, user)
	}
	return users, nil
}

func (r *repo) CreateChat(req CreateChatReq) (*chat.Chat, error) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	id, err := r.createChat(tx, req.Name, req.OwnerID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return nil, err
	}

	req.UsersIDs = append(req.UsersIDs, req.OwnerID)
	err = r.addChatUsers(tx, id, req.UsersIDs)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	users := make([]*user.User, 0, len(req.UsersIDs))
	for _, id := range req.UsersIDs {
		users = append(users, user.New(id))
	}

	return chat.New(id, req.Name, user.New(req.OwnerID), users), nil
}

func (r *repo) createChat(tx *sql.Tx, name string, ownerID int) (int, error) {
	query := `INSERT INTO chats (name, owner) VALUES ($1, $2) RETURNING id`
	res := tx.QueryRow(query, name, ownerID)

	var id int
	if err := res.Scan(&id); err != nil {
		log.Error(errors.WithStack(err))
		return 0, common.ErrInternal
	}
	
	return id, nil
}

func (r *repo) addChatUsers(tx *sql.Tx, chatID int, usersIDs []int) error {
	var b bytes.Buffer
	b.WriteString(`INSERT INTO chats_users (chat_id, user_id) VALUES `)
	rawChatID := strconv.Itoa(chatID)
	for i, userID := range usersIDs {
		if (i != 0){
			b.WriteString(",")
		}
		b.WriteString("(")
		b.WriteString(rawChatID)
		b.WriteString(",")
		b.WriteString(strconv.Itoa(userID))
		b.WriteString(")")
	}

	query := b.String()
	_, err := tx.Exec(query)
	if err == sql.ErrTxDone {
		return common.ErrInternal
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}
