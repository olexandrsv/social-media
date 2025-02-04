package repository

import (
	"bytes"
	"database/sql"
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
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
	UsersFullNames([]int) ([]*user.User, error)
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

func (r *repo) SubscriptionExists(userID, followedID int) (bool, error){
	query := `select count(user_id) from followers where user_id=$1 and follower_id=$2`
	row := r.db.QueryRow(query, followedID, userID)

	var count int
	if err := row.Scan(&count); err != nil{
		log.Error(errors.WithStack(err))
		return false, common.ErrInternal
	}

	if count == 0{
		return false, nil
	}
	return true, nil
}

func (r *repo) Subscribe(userID, followedID int) error {
	query := `insert into followers (read, user_id, follower_id) values ($1, $2, $3)`
	_, err := r.db.Exec(query, 0, followedID, userID)
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

func (r *repo) UsersFullNames(ids []int) ([]*user.User, error){
	var b bytes.Buffer
	var caseBuffer bytes.Buffer
	for i, id := range ids {
		if i != 0 {
			b.WriteString(", ")
		}
		v := strconv.Itoa(id)
		b.WriteString(v)
		caseBuffer.WriteString(" when ")
		caseBuffer.WriteString(v)
		caseBuffer.WriteString(" then ")
		caseBuffer.WriteString(strconv.Itoa(i+1))
	}
	caseBuffer.WriteString(" end;")

	query := fmt.Sprintf("select id, first_name, second_name from users where id in (%s) order by case id %s", 
		b.String(), caseBuffer.String())

	rows, err := r.db.Query(query)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, common.ErrInternal
	}

	users := make([]*user.User, 0, len(ids))
	for rows.Next() {
		var model UserModel
		if err := rows.Scan(&model.ID, &model.Name, &model.Surname); err != nil {
			return nil, err
		}
		u := user.New(model.ID, "", user.WithName(model.Name), user.WithSurname(model.Surname))
		users = append(users, u)
	}

	return users, nil
}