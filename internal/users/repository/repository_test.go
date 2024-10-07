package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"social-media/internal/users/domain/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type QueryType int

const (
	QueryRows QueryType = iota
	Exec
)

type Query struct {
	sql  string
	t    QueryType
	args []driver.Value
}

type QueryResult struct {
	header []string
	data   [][]driver.Value
	err    error
}

func newSqlMock(query Query, res QueryResult) (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	switch {
	case query.t == QueryRows:
		if res.err != nil {
			mock.ExpectQuery(query.sql).WithArgs(query.args...).WillReturnError(res.err)
			break
		}
		rows := mock.NewRows(res.header).AddRows(res.data...)
		mock.ExpectQuery(query.sql).WithArgs(query.args...).WillReturnRows(rows)
	case query.t == Exec:
		if res.err != nil {
			mock.ExpectExec(query.sql).WithArgs(query.args...).WillReturnError(res.err)
			break
		}
		mock.ExpectExec(query.sql).WithArgs(query.args...).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	return db, mock, nil
}

func TestCreateUser(t *testing.T) {
	u := user.NewUser(2, "bob", user.WithName("Bob"),
		user.WithSurname("Smith"), user.WithPassword("12her45j"))

	userModel := NewUserModel(u.Login(), u.Name(), u.Surname(), u.Password())

	query := Query{
		sql:  "insert into users (.+) values (.+) returning id",
		t:    QueryRows,
		args: []driver.Value{u.Login(), u.Name(), u.Surname(), u.Password()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes    QueryResult
		createdUser *user.User
	}{
		{
			queryRes: QueryResult{
				header: []string{"id"},
				data: [][]driver.Value{
					{2},
				},
				err: nil,
			},
			createdUser: u,
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			createdUser: nil,
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}

		r := repo{db}
		createdUser, err := r.CreateUser(userModel)
		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err != nil{
			return
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if u.Login() != createdUser.Login() || u.Name() != createdUser.Name() ||
			u.Surname() != createdUser.Surname() || u.Password() != createdUser.Password() {
			t.Errorf("Incorrect user: %+v, expected %+v", createdUser, u)
		}
	}
}

func TestGetCredentials(t *testing.T) {
	user := user.NewUser(12, "bob",
		user.WithPassword("iYX1Jc1jT+i4Fb2hz7us4/H2w3DL6wOrYNWXJqKJvovcbv7orTULJO843frjFt+B"))

	query := Query{
		sql:  `select (.+) from users where login=?`,
		t:    QueryRows,
		args: []driver.Value{user.Login()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes    QueryResult
		resID       int
		resPassword string
	}{
		{
			queryRes: QueryResult{
				header: []string{"id", "password"},
				data: [][]driver.Value{
					{user.ID(), user.Password()},
				},
				err: nil,
			},
			resID:       user.ID(),
			resPassword: user.Password(),
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			resID:       0,
			resPassword: "",
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}

		r := repo{db}
		id, password, err := r.GetCredentials(user.Login())

		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if id != d.resID || password != d.resPassword {
			t.Errorf("error: id: %d, expected: %d; password: %s expected: %s",
				id, user.ID(), password, user.Password())
		}
	}
}

func TestGetUser(t *testing.T) {
	u := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"),
		user.WithBio("student"), user.WithInterests("play chess"))

	query := Query{
		sql:  `select (.+) from users where login=?`,
		t:    QueryRows,
		args: []driver.Value{u.Login()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes QueryResult
		resUser  *user.User
	}{
		{
			queryRes: QueryResult{
				header: []string{"id", "login", "first_name", "second_name", "bio", "interests"},
				data: [][]driver.Value{
					{u.ID(), u.Login(), u.Name(), u.Surname(), u.Bio(), u.Interests()},
				},
				err: nil,
			},
			resUser: u,
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			resUser: nil,
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			t.Fatalf(err.Error())
		}

		r := repo{db}
		user, err := r.GetUser(u.Login())

		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if user == d.resUser {
			return
		}

		if user.Name() != d.resUser.Name() ||
			user.Surname() != d.resUser.Surname() ||
			user.Bio() != d.resUser.Bio() ||
			user.Interests() != d.resUser.Interests() {
			t.Errorf("error: user: %v, expected: %v", user, d.resUser)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	u := user.NewUser(1, "bob", user.WithName("Bob"), user.WithSurname("Smith"),
		user.WithBio("student"), user.WithInterests("play chess"))

	query := Query{
		sql:  `update users set (.+) where id=?`,
		t:    Exec,
		args: []driver.Value{u.Name(), u.Surname(), u.Bio(), u.Interests(), u.ID()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes QueryResult
	}{
		{
			queryRes: QueryResult{
				err: nil,
			},
		},
		{
			queryRes: QueryResult{
				err: e,
			},
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}
		r := repo{db}
		err = r.UpdateUser(u)

		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	}
}

func TestUserExists(t *testing.T) {
	u := user.NewUser(1, "bob")

	query := Query{
		sql:  `select (.+) from users where login=?`,
		t:    QueryRows,
		args: []driver.Value{u.Login()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes QueryResult
		exists   bool
	}{
		{
			queryRes: QueryResult{
				header: []string{"count"},
				data: [][]driver.Value{
					{0},
				},
				err: nil,
			},
			exists: false,
		},
		{
			queryRes: QueryResult{
				header: []string{"count"},
				data: [][]driver.Value{
					{1},
				},
				err: nil,
			},
			exists: true,
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			exists: false,
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}

		r := repo{db}
		exists, err := r.UserExists(u.Login())

		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if exists != d.exists {
			t.Errorf("error: exists is %t, expected: %t", exists, d.exists)
		}
	}
}

func TestGetLoginsByInfo(t *testing.T) {
	info := "book"
	query := Query{
		sql:  `select (.+) from users`, //  where interests like $1 or bio like $1
		t:    QueryRows,
		args: []driver.Value{"%" + info + "%"},
	}
	e := errors.New("postgres error")
	scanErr := errors.New(`sql: Scan error on column index 0, name "login": converting NULL to string is unsupported`)

	data := []struct {
		queryRes  QueryResult
		resLogins []string
		resErr    error
	}{
		{
			queryRes: QueryResult{
				header: []string{"login"},
				data: [][]driver.Value{
					{"bob"},
					{"ben"},
					{"bill"},
				},
				err: nil,
			},
			resLogins: []string{"bob", "ben", "bill"},
			resErr:    nil,
		},
		{
			queryRes: QueryResult{
				header: []string{"login"},
				data: [][]driver.Value{
					{"smith"},
				},
				err: nil,
			},
			resLogins: []string{"smith"},
			resErr:    nil,
		},
		{
			queryRes: QueryResult{
				header: []string{"login"},
				data: [][]driver.Value{
					{nil},
				},
				err: nil,
			},
			resLogins: nil,
			resErr:    scanErr,
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			resLogins: nil,
			resErr:    e,
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}

		r := repo{db}
		logins, err := r.GetLoginsByInfo(info)

		if !errors.Is(err, d.resErr) && err.Error() != d.resErr.Error() {
			t.Errorf("error: '%s' expected: %s", err, d.resErr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		for i, v := range logins {
			if v != d.resLogins[i] {
				t.Errorf("incorrect logins %v, expected: %v", logins, d.resLogins)
			}
		}
	}
}

func TestSubsribe(t *testing.T) {
	bob := user.NewUser(1, "bob")
	ben := user.NewUser(2, "ben")

	query := Query{
		sql:  `insert into followers`, // (.+) values ((select id from users where login=?), ?)`,
		t:    Exec,
		args: []driver.Value{ben.Login(), bob.ID()},
	}
	e := errors.New("postgres error")

	data := []struct {
		queryRes QueryResult
	}{
		{
			queryRes: QueryResult{
				err: nil,
			},
		},
		{
			queryRes: QueryResult{
				err: e,
			},
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}
		r := repo{db}
		err = r.Subscribe(bob.ID(), ben.Login())

		if !errors.Is(err, d.queryRes.err) {
			t.Errorf("error: '%s' expected: %s", err, d.queryRes.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetFollowedLogins(t *testing.T) {
	u := user.NewUser(1, "bob")
	query := Query{
		sql:  `select (.+) from users join followers on users.id = followers.user_id and followers.follower_id=?`,
		t:    QueryRows,
		args: []driver.Value{u.ID()},
	}
	e := errors.New("postgres error")
	scanErr := errors.New(`sql: Scan error on column index 0, name "login": converting NULL to string is unsupported`)

	data := []struct {
		queryRes  QueryResult
		resLogins []string
		resErr    error
	}{
		{
			queryRes: QueryResult{
				header: []string{"login"},
				data: [][]driver.Value{
					{"bob"},
					{"ben"},
					{"bill"},
				},
				err: nil,
			},
			resLogins: []string{"bob", "ben", "bill"},
			resErr:    nil,
		},
		{
			queryRes: QueryResult{
				header: []string{"login"},
				data: [][]driver.Value{
					{nil},
				},
				err: nil,
			},
			resLogins: nil,
			resErr:    scanErr,
		},
		{
			queryRes: QueryResult{
				err: e,
			},
			resLogins: nil,
			resErr:    e,
		},
	}

	for _, d := range data {
		db, mock, err := newSqlMock(query, d.queryRes)
		if err != nil {
			log.Fatal(err)
		}

		r := repo{db}
		logins, err := r.GetFollowedLogins(u.ID())

		if !errors.Is(err, d.resErr) && err.Error() != d.resErr.Error() {
			t.Errorf("error: '%s' expected: %s", err, d.resErr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if len(logins) != len(d.resLogins) {
			t.Errorf("incorrect logins %+q, expected: %+q", logins, d.resLogins)
			return
		}

		for i, v := range logins {
			if v != d.resLogins[i] {
				t.Errorf("incorrect logins %+q, expected: %+q", logins, d.resLogins)
			}
		}
	}

}
