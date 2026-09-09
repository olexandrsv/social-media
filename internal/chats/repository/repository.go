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
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository interface {
	UserChats(int) ([]*chat.Chat, error)
	Chat(int) (*chat.Chat, error)
	CreateChat(CreateChatReq) (*chat.Chat, error)
	UpdateChat(*chat.Chat) error
	ChatOwner(int) (int, error)
	DeleteChat(int) error
	UpdateReadMessages(int, int, string) error
	LastReadMessages(int) ([]*chatmessage.ChatMessage, error)
}

type repo struct {
	db *sql.DB
}

func New() (Repository, error) {
	password, err := common.ReadSecret("chats-database-secret")
	if err != nil {
		return nil, err
	}
	user := config.App.Chats.DB.User
	host := config.App.Chats.DB.Host
	port := config.App.Chats.DB.Port
	name := config.App.Chats.DB.Name
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error(err)
		panic("Unable to connect to database")
	}
	return &repo{
		db: db,
	}, nil
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

func (r *repo) Chat(id int) (*chat.Chat, error) {
	c, err := r.getChat(id)
	if err != nil {
		return nil, err
	}
	users, err := r.getChatUsers(id)
	if err != nil {
		return nil, err
	}

	return chat.New(c.ID(), c.Name(), c.Owner(), users), nil
}

func (r *repo) getChat(id int) (*chat.Chat, error) {
	query := `SELECT name, owner FROM chats WHERE id=$1`
	row := r.db.QueryRow(query, id)

	var ownerID int
	var name string
	if err := row.Scan(&name, &ownerID); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	owner := user.New(ownerID)
	chat := chat.New(id, name, owner, nil)
	return chat, nil
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
	b.WriteString(`INSERT INTO chats_users (chat_id, user_id, last_read_message) VALUES `)
	rawChatID := strconv.Itoa(chatID)
	for i, userID := range usersIDs {
		if i != 0 {
			b.WriteString(",")
		}
		b.WriteString("(")
		b.WriteString(rawChatID)
		b.WriteString(",")
		b.WriteString(strconv.Itoa(userID))
		b.WriteString(",")
		b.WriteString("''")
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

func (r *repo) UpdateChat(c *chat.Chat) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	err = r.updateChat(tx, c.ID(), c.Name())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return err
	}

	err = r.deleteChatUsers(tx, c.ID())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return err
	}

	ids := slice.MustConvert(c.Users(), func(u *user.User) int {
		return u.ID()
	})
	ids = append(ids, c.Owner().ID())
	err = r.addChatUsers(tx, c.ID(), ids)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) updateChat(tx *sql.Tx, chatID int, name string) error {
	query := `UPDATE chats SET name=$1 WHERE id=$2`

	_, err := tx.Exec(query, name, chatID)
	if err == sql.ErrTxDone {
		return common.ErrInternal
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) deleteChatUsers(tx *sql.Tx, chatID int) error {
	query := `DELETE FROM chats_users WHERE chat_id=$1`

	_, err := tx.Exec(query, chatID)
	if err == sql.ErrTxDone {
		return common.ErrInternal
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) ChatOwner(chatID int) (int, error) {
	query := `SELECT owner FROM chats WHERE id=$1`

	row := r.db.QueryRow(query, chatID)

	var ownerID int
	if err := row.Scan(&ownerID); err != nil {
		log.Error(errors.WithStack(err))
		return 0, common.ErrInternal
	}

	return ownerID, nil
}

func (r *repo) DeleteChat(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	err = r.deleteChat(tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return err
	}

	err = r.deleteChatUsers(tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(errors.WithStack(err))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) deleteChat(tx *sql.Tx, id int) error {
	query := `DELETE FROM chats WHERE id=$1`

	_, err := tx.Exec(query, id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) UpdateReadMessages(chatID, userID int, messageID string) error {
	query := `UPDATE chats_users SET last_read_message=$1 WHERE chat_id=$2 AND user_id=$3`
	_, err := r.db.Exec(query, messageID, chatID, userID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) LastReadMessages(userID int) ([]*chatmessage.ChatMessage, error) {
	query := `SELECT chat_id, last_read_message FROM chats_users WHERE user_id=$1`
	rows, err := r.db.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var messages []*chatmessage.ChatMessage
	for rows.Next() {
		var m MessageModel
		if err := rows.Scan(&m.ChatID, &m.ID); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		messages = append(messages, chatmessage.New(m.ID, 0, "", nil, nil, m.ChatID))
	}

	return messages, nil
}
