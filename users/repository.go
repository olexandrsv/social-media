package users

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveUser(*User) (int, error)
	GetCredentials(string) (int, string, error)
	GetUser(string) (*User, error)
	UpdateUser(*User) error
	GetLoginsByInfo(string) ([]string, error)
	Subscribe(int, int) error
	GetIdByLogin(string) (int, error)
	GetFollowedLogins(int) ([]string, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) SaveUser(u *User) (int, error) {
	sql := `insert into users (login, first_name, second_name, password) 
		values ($1, $2, $3, $4) returning id`
	var id int
	err := r.db.QueryRow(context.Background(), sql, u.Login, u.FirstName, u.SecondName, u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) GetCredentials(login string) (int, string, error) {
	sql := `select id, password from users where login=$1`
	var id int
	var encodedPassw string

	err := r.db.QueryRow(context.Background(), sql, login).Scan(&id, &encodedPassw)
	if err != nil {
		return 0, "", err
	}
	return id, encodedPassw, nil
}

func (r *repo) GetUser(login string) (*User, error) {
	sql := `select first_name, second_name, bio, interests from users where login=$1`
	user := NewUser(login)
	err := r.db.QueryRow(context.Background(), sql, login).Scan(&user.FirstName,
		&user.SecondName, &user.Bio, &user.Interests)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repo) UpdateUser(u *User) error {
	sql := `update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5`
	_, err := r.db.Exec(context.Background(), sql, u.FirstName, u.SecondName, u.Bio, u.Interests, u.ID)
	return err
}

func (r *repo) GetLoginsByInfo(info string) ([]string, error) {
	sql := `select login from users where interests like $1 or bio like $1`
	rows, err := r.db.Query(context.Background(), sql, "%"+info+"%")
	if err != nil {
		return nil, err
	}

	var logins []string
	for rows.Next() {
		var login string
		err = rows.Scan(&login)
		if err != nil {
			continue
		}
		logins = append(logins, login)
	}
	return logins, nil
}

func (r *repo) Subscribe(userID, followedID int) error{
	sql := `insert into followers (user_id, follower_id) values ($1, $2)`
	_, err := r.db.Exec(context.Background(), sql, followedID, userID)
	return err
}

func (r *repo) GetIdByLogin(login string) (int, error) {
	sql := `select id from users where login=$1`
	var id int
	err := r.db.QueryRow(context.Background(), sql, login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) GetFollowedLogins(id int) ([]string, error){
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
			log.Println(err)
			continue
		}
		logins = append(logins, login)
	}
	return logins, nil
}

