package users

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveUser(*User) (int, error)
	GetCredentials(string) (int, string, error)
	GetUser(string) (*User, error)
	UpdateUser(*User) error
	UserExists(string) (bool, error)
	GetLoginsByInfo(string) ([]string, error)
	Subscribe(int, string) error
	GetFollowedLogins(int) ([]string, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(url string) Repository {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Println(err)
		panic("Unable to connect to database")
	}
	return &repo{
		db: pool,
	}
}

func (r *repo) SaveUser(u *User) (int, error) {
	sql := `insert into users (login, first_name, second_name, password) 
		values ($1, $2, $3, $4) returning id`
	var id int
	err := r.db.QueryRow(context.Background(), sql, u.Login, u.FirstName, u.SecondName, u.Password).Scan(&id)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return id, nil
}

func (r *repo) GetCredentials(login string) (int, string, error) {
	sql := `select id, password from users where login=$1`
	var id int
	var encodedPassw string

	err := r.db.QueryRow(context.Background(), sql, login).Scan(&id, &encodedPassw)
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	return id, encodedPassw, nil
}

func (r *repo) GetUser(login string) (*User, error) {
	sql := `select first_name, second_name, bio, interests from users where login=$1`
	user := NewUser(login)
	err := r.db.QueryRow(context.Background(), sql, login).Scan(&user.FirstName,
		&user.SecondName, &user.Bio, &user.Interests)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (r *repo) UpdateUser(u *User) error {
	sql := `update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5`
	_, err := r.db.Exec(context.Background(), sql, u.FirstName, u.SecondName, u.Bio, u.Interests, u.ID)
	return errors.WithStack(err)
}

func (r *repo) UserExists(login string) (bool, error) {
	sql := `select count(*) from users where login=$1`
	var count int
	err := r.db.QueryRow(context.Background(), sql, login).Scan(&count)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repo) GetLoginsByInfo(info string) ([]string, error) {
	sql := `select login from users where interests like $1 or bio like $1`
	rows, err := r.db.Query(context.Background(), sql, "%"+info+"%")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var logins []string
	for rows.Next() {
		var login string
		err = rows.Scan(&login)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		logins = append(logins, login)
	}
	return logins, nil
}

func (r *repo) Subscribe(userID int, followedLogin string) error {
	sql := `insert into followers (user_id, follower_id) values ((select id from users where login=$1), $2)`
	_, err := r.db.Exec(context.Background(), sql, followedLogin, userID)
	return errors.WithStack(err)
}

func (r *repo) GetFollowedLogins(id int) ([]string, error) {
	sql := `select login from users join followers on users.id = followers.user_id and followers.follower_id=$1`
	rows, err := r.db.Query(context.Background(), sql, id)
	if err != nil {
		return nil, err
	}

	var logins []string
	for rows.Next() {
		var login string
		err = rows.Scan(&login)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		logins = append(logins, login)
	}
	return logins, nil
}
