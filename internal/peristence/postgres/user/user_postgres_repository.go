package chat

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/config"
	"ru.nklimkin/petmsngr/internal/domain/user"

	_ "github.com/lib/pq"
)

const DATABASE_CONNECTION_TEMPLATE = "postgres://%s:%s@%s/%s?sslmode=disable"

type UserPostgresRepository struct {
	db *sqlx.DB
}

func New(datasource config.Datasource) (*UserPostgresRepository, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			DATABASE_CONNECTION_TEMPLATE, datasource.Username, datasource.Password, datasource.Host, datasource.DatabaseName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS messanger_user (
		id BIGINT PRIMARY KEY,
		login VARCHAR(255),
		created TIMESTAMP);
	`)

	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to create table - messanger_user, error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("can't create table - messanger_user, error: %w", err)
	}

	stmt, err = db.Prepare(`CREATE SEQUENCE IF NOT EXISTS message_user_id_seq START 1;`)

	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to create sequence - message_user_id_seq, error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("can't create sequence - message_user_id_seq, error: %w", err)
	}

	return &UserPostgresRepository{db}, nil
}

func (rep *UserPostgresRepository) GetById(id user.UserId) (*user.User, error) {
	stmt, err := rep.db.Prepare("SELECT * FROM messanger_user WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to get user by id: %w", err)
	}

	var persistedId int64
	var persistedLogin string
	var persistedCreated time.Time

	err = stmt.QueryRow(id.Value).Scan(&persistedId, &persistedLogin, &persistedCreated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("there is no user with id = [%d]", id.Value)
	}
	if err != nil {
		return nil, fmt.Errorf("can't get user with id [%d]: %w", id.Value, err)
	}

	userId := user.UserId{Value: persistedId}

	return &user.User{Id: userId, Login: persistedLogin}, nil
}

func (rep *UserPostgresRepository) Save(user user.User) (*user.User, error) {
	stmt, err := rep.db.Prepare("INSERT INTO messanger_user(id, login, created) VALUES($1, $2, $3)")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare query to save user: %w", err)
	}
	_, err = stmt.Exec(user.Id.Value, user.Login, time.Now())
	if err != nil {
		return nil, fmt.Errorf("can't save user: %w", err)
	}
	return &user, nil
}

func (rep *UserPostgresRepository) Generate() (*user.UserId, error) {
	stmt, err := rep.db.Prepare("SELECT nextval('message_user_id_seq')")
	if err != nil {
		return nil, fmt.Errorf("can't builde prepate statement to get new user id, error: %w", err)
	}

	var id int64

	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("can't get new user id, error: %w", err)
	}

	return &user.UserId{Value: id}, nil
}

func (rep *UserPostgresRepository) Exists(id user.UserId) (bool, error) {
	stmt, err := rep.db.Prepare("SELECT id FROM messanger_user WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("can't build prepate statement to check user existence, error: %w", err)
	}

	var persistedId int64

	err = stmt.QueryRow(id.Value).Scan(&persistedId)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't check user existence, error: %w", err)
	}

	return true, nil
}
