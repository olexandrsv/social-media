package repository

import (
	"bytes"
	"database/sql"
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	"social-media/internal/users/domain/post"
	"social-media/internal/users/domain/user"
	"strconv"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type Repository interface {
	CreateUser(UserModel) (*user.User, error)
	GetCredentials(string) (int, string, error)
	GetUser(int) (*user.User, error)
	UpdateUser(*user.User) error
	UserExists(string) (bool, error)
	GetUsersByInfo(string) ([]*user.User, error)
	SubscriptionExists(int, int) (bool, error)
	Subscribe(int, int) error
	GetFollowedUsers(int) ([]*user.User, error)
	UsersInfo([]int) ([]*user.User, error)
	UpdateReadPosts(int, int, string) error
	LastReadPosts(int) ([]*post.Post, error)
	GetFollowers(int) ([]*user.User, error)
}

type repo struct {
	db *sql.DB
}

func New() Repository {
	user := config.App.Users.DB.User
	password := config.App.Users.DB.Password
	host := config.App.Users.DB.Host
	port := config.App.Users.DB.Port
	name := config.App.Users.DB.Name
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	log.Infof("URL: %s", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error(err)
		panic("Unable to connect to database")
	}
	return &repo{
		db: db,
	}
}

func (r *repo) CreateUser(userModel UserModel) (*user.User, error) {
	sql := `insert into users (login, first_name, second_name, password) 
		values ($1, $2, $3, $4) returning id`
	var id int
	err := r.db.QueryRow(sql, userModel.Login, userModel.Name, userModel.Surname, userModel.Password).Scan(&id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	return user.New(id, userModel.Login, user.WithName(userModel.Name),
		user.WithSurname(userModel.Surname), user.WithPassword(userModel.Password)), nil
}

func (r *repo) GetCredentials(login string) (int, string, error) {
	query := `select id, password from users where login=$1`
	var id int
	var encodedPassw string

	err := r.db.QueryRow(query, login).Scan(&id, &encodedPassw)
	if err == sql.ErrNoRows {
		return 0, "", common.ErrWrongCredentials
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return 0, "", common.ErrInternal
	}
	return id, encodedPassw, nil
}

func (r *repo) GetUser(id int) (*user.User, error) {
	query := `select login, first_name, second_name, bio, interests from users where id=$1`
	userModel := UserModel{
		ID: id,
	}
	err := r.db.QueryRow(query, id).Scan(&userModel.Login, &userModel.Name,
		&userModel.Surname, &userModel.Bio, &userModel.Interests)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	user := user.New(userModel.ID, userModel.Login, user.WithName(userModel.Name),
		user.WithSurname(userModel.Surname), user.WithBio(userModel.Bio), user.WithInterests(userModel.Interests))
	return user, nil
}

func (r *repo) UpdateUser(u *user.User) error {
	query := `update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5`
	_, err := r.db.Exec(query, u.Name(), u.Surname(), u.Bio(), u.Interests(), u.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) UserExists(login string) (bool, error) {
	query := `select count(id) from users where login=$1`
	var count int
	err := r.db.QueryRow(query, login).Scan(&count)
	if err != nil {
		log.Error(errors.WithStack(err))
		return false, common.ErrInternal
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repo) GetUsersByInfo(info string) ([]*user.User, error) {
	query := `select id, login from users where login like $1 or interests like $1 or bio like $1 limit 7`
	rows, err := r.db.Query(query, "%"+info+"%")
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var users []*user.User
	for rows.Next() {
		var userModel UserModel
		err = rows.Scan(&userModel.ID, &userModel.Login)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		users = append(users, user.New(userModel.ID, userModel.Login))
	}

	return users, nil
}

func (r *repo) SubscriptionExists(userID, followedID int) (bool, error) {
	query := `select count(user_id) from followers where user_id=$1 and follower_id=$2`
	row := r.db.QueryRow(query, followedID, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		log.Error(errors.WithStack(err))
		return false, common.ErrInternal
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repo) Subscribe(userID, followedID int) error {
	query := `insert into followers (last_read_post, user_id, follower_id) values ($1, $2, $3)`
	_, err := r.db.Exec(query, "", followedID, userID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) GetFollowedUsers(id int) ([]*user.User, error) {
	query := `select id, login from users join followers on users.id = followers.user_id and followers.follower_id=$1`
	rows, err := r.db.Query(query, id)
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
		var login string
		err = rows.Scan(&id, &login)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		users = append(users, user.New(id, login))
	}
	return users, nil
}

func (r *repo) UsersInfo(ids []int) ([]*user.User, error) {
	if len(ids) == 0 {
		return []*user.User{}, nil
	}
	var b bytes.Buffer
	for i, id := range ids {
		if i != 0 {
			b.WriteString(", ")
		}
		v := strconv.Itoa(id)
		b.WriteString("(")
		b.WriteString(v)
		b.WriteString(",")
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(")")
	}
	b.WriteString(";")

	createTable := `CREATE TEMP TABLE users_data (
		user_id INTEGER,
		idx INTEGER
	);`
	insert := `INSERT INTO users_data (user_id, idx) VALUES ` + b.String()
	get := `SELECT id, login, first_name, second_name FROM users_data LEFT JOIN users ON users_data.user_id = users.id ORDER BY users_data.idx;`
	dropTable := `DROP TABLE users_data;`

	query := createTable + insert + get + dropTable

	rows, err := r.db.Query(query)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	users := make([]*user.User, 0, len(ids))
	for rows.Next() {
		var model UserModel
		if err := rows.Scan(&model.ID, &model.Login, &model.Name, &model.Surname); err != nil {
			return nil, err
		}
		u := user.New(model.ID, model.Login, user.WithName(model.Name), user.WithSurname(model.Surname))
		users = append(users, u)
	}

	return users, nil
}

func (r *repo) UpdateReadPosts(ownerID, userID int, messageID string) error {
	query := `UPDATE followers SET last_read_post=$1 WHERE user_id=$2 AND follower_id=$3`
	_, err := r.db.Exec(query, messageID, ownerID, userID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) LastReadPosts(userID int) ([]*post.Post, error) {
	query := `SELECT user_id, last_read_post FROM followers WHERE follower_id=$1`
	rows, err := r.db.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}

	var posts []*post.Post
	for rows.Next() {
		var p PostModel
		if err := rows.Scan(&p.UserID, &p.ID); err != nil {
			return nil, common.ErrInternal
		}
		posts = append(posts, post.NewPost(p.ID, p.UserID))
	}
	return posts, nil
}

func (r *repo) GetFollowers(userID int) ([]*user.User, error) {
	query := `SELECT follower_id FROM followers WHERE user_id=$1`
	rows, err := r.db.Query(query, userID)
	if err == sql.ErrNoRows {
		return []*user.User{}, nil
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var users []*user.User
	for rows.Next() {
		var u UserModel
		if err = rows.Scan(&u.ID); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		users = append(users, user.New(u.ID, ""))
	}

	return users, nil
}
