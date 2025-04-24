package chat

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type UserPostgresRepository struct {
	db *sqlx.DB
}

func New() (*UserPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=messanger sslmode=disable password=postgres host=localhost")
	if err != nil {
		panic("Can't connect to db")
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS messanger_user (
		id BIGINT PRIMARY KEY,
		login VARCHAR(255),
		created TIMESTAMP);
		`)

	if err != nil {
		panic("Can't create messanger_user table")
	}

	_, err = stmt.Exec()
	if err != nil {
		panic("Can't execute query to create messanger_user table")
	}

	return &UserPostgresRepository{db}, nil
}

func (rep *UserPostgresRepository) GetById(id user.UserId) user.User {
	stmt, err := rep.db.Prepare("SELECT * FROM messanger_user WHERE id = $1")
	if err != nil {
		panic("Can't preapre query to get messagenger_user")
	}

	var persistedId int64
	var persistedLogin string
	var persistedCreated time.Time

	err = stmt.QueryRow(id.Value).Scan(&persistedId, &persistedLogin, &persistedCreated)
	if errors.Is(err, sql.ErrNoRows) {
		panic("No messanger user")
	}
	if err != nil {
		panic("Can't get messanger user")
	}

	userId := user.UserId{Value: persistedId}

	return user.User{Id: userId, Login: persistedLogin}
}

func (rep *UserPostgresRepository) Save(user user.User) {
	stmt, err := rep.db.Prepare("INSERT INTO messanger_user(id, login, creater) VALUES($1, $2, $3)")
	if err != nil {
		panic("Can't insert user to db")
	}
	_, err = stmt.Exec(user.Id.Value, user.Login, time.Now())
	if err != nil {
		panic("Can't insert user to db")
	}
}
