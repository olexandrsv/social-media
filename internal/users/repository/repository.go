package repository

import (
	"database/sql"
	"fmt"
	"log"
	"social-media/internal/common/app/config"
	"social-media/internal/users/domain/user"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type Repository interface {
	CreateUser(UserModel) (*user.User, error)
	GetCredentials(string) (int, string, error)
	GetUser(string) (*user.User, error)
	UpdateUser(*user.User) error
	UserExists(string) (bool, error)
	GetLoginsByInfo(string) ([]string, error)
	Subscribe(int, string) error
	GetFollowedLogins(int) ([]string, error)
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
		log.Println(err)
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
		return nil, errors.WithStack(err)
	}
	return user.NewUser(id, userModel.Login, user.WithName(userModel.Name),
		user.WithSurname(userModel.Surname), user.WithPassword(userModel.Password)), nil
}

func (r *repo) GetCredentials(login string) (int, string, error) {
	sql := `select id, password from users where login=$1`
	var id int
	var encodedPassw string

	err := r.db.QueryRow(sql, login).Scan(&id, &encodedPassw)
	if err != nil {
		return 0, "", err
	}
	return id, encodedPassw, nil
}

func (r *repo) GetUser(login string) (*user.User, error) {
	sql := `select id, login, first_name, second_name, bio, interests from users where login=$1`
	userModel := UserModel{}
	err := r.db.QueryRow(sql, login).Scan(&userModel.ID, &userModel.Login, &userModel.Name,
		&userModel.Surname, &userModel.Bio, &userModel.Interests)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user := user.NewUser(userModel.ID, userModel.Login, user.WithName(userModel.Name),
		user.WithSurname(userModel.Surname), user.WithBio(userModel.Bio), user.WithInterests(userModel.Interests))
	return user, nil
}

func (r *repo) UpdateUser(u *user.User) error {
	sql := `update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5`
	_, err := r.db.Exec(sql, u.Name(), u.Surname(), u.Bio(), u.Interests(), u.ID())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *repo) UserExists(login string) (bool, error) {
	sql := `select count(*) from users where login=$1`
	var count int
	err := r.db.QueryRow(sql, login).Scan(&count)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repo) GetLoginsByInfo(info string) ([]string, error) {
	query := `select login from users where interests like $1 or bio like $1`
	rows, err := r.db.Query(query, "%"+info+"%")
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
	_, err := r.db.Exec(sql, followedLogin, userID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *repo) GetFollowedLogins(id int) ([]string, error) {
	sql := `select login from users join followers on users.id = followers.user_id and followers.follower_id=$1`
	rows, err := r.db.Query(sql, id)
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
