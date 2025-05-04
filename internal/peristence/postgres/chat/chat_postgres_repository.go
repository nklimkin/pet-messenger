package chat

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/config"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"

	_ "github.com/lib/pq"
)

const DATABASE_CONNECTION_TEMPLATE = "postgres://%s:%s@%s/%s?sslmode=disable"

type ChatPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(datasource config.Datasource) (*ChatPostgresRepository, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			DATABASE_CONNECTION_TEMPLATE, datasource.Username, datasource.Password, datasource.Host, datasource.DatabaseName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS chats (
		id BIGINT PRIMARY KEY,
		first_user_id BIGINT REFERENCES messanger_user(id),
		second_user_id BIGINT REFERENCES messanger_user(id),
		created TIMESTAMP);
	`)

	if err != nil {
		return nil, fmt.Errorf("can't build statement to create table - chats, error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("can't create table - chats, error: %w", err)
	}

	stmt, err = db.Prepare(`CREATE SEQUENCE IF NOT EXISTS chat_id_seq START 1;`)

	if err != nil {
		return nil, fmt.Errorf("can't build statement to create sequence - chat_id_seq, error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("can't create sequence - chat_id_seq, error: %w", err)
	}

	return &ChatPostgresRepository{db}, nil
}

func (rep *ChatPostgresRepository) GetById(id chat.ChatId) (*chat.Chat, error) {
	stmt, err := rep.db.Prepare("SELECT * FROM chats WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to get chat by id, error: %w", err)
	}

	var persistedId int64
	var persistedFirstUserId int64
	var persistedSecondUserId int64
	var persistedCreated time.Time

	err = stmt.QueryRow(id.Value).Scan(&persistedId, &persistedFirstUserId, &persistedSecondUserId, &persistedCreated)
	if err != nil {
		return nil, fmt.Errorf("can't build chat from result set: %w", err)
	}

	return buildChat(persistedId, persistedFirstUserId, persistedSecondUserId, persistedCreated), nil
}

func (rep *ChatPostgresRepository) GetByUserId(userId user.UserId) ([]*chat.Chat, error) {
	stmt, err := rep.db.Prepare("SELECT * FROM chats WHERE first_user_id = $1 OR second_user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statment to get chats by user id, error: %w", err)
	}

	rows, err := stmt.Query(userId.Value)
	if err != nil {
		return nil, fmt.Errorf("can't get chats by user id = [%d], error: %w", userId.Value, err)
	}

	chats := make([]*chat.Chat, 1)

	for rows.Next() {
		var persistedId int64
		var persistedFirstUserId int64
		var persistedSecondUserId int64
		var persistedCreated time.Time

		err = rows.Scan(&persistedId, &persistedFirstUserId, &persistedSecondUserId, &persistedCreated)
		if err != nil {
			return nil, fmt.Errorf("can't build chat from result set, error: %w", err)
		}

		chats = append(chats, buildChat(persistedId, persistedFirstUserId, persistedSecondUserId, persistedCreated))
	}

	return chats, nil
}

func buildChat(
	persistedId int64,
	persistedFirstUserId int64,
	persistedSecondUserId int64,
	persistedCreated time.Time,
) *chat.Chat {
	chatId := chat.ChatId{Value: persistedId}
	firstUserId := user.UserId{Value: persistedFirstUserId}
	secondUserId := user.UserId{Value: persistedSecondUserId}
	return &chat.Chat{Id: chatId, FirstUser: firstUserId, SecondUser: secondUserId, Created: persistedCreated}
}

func (rep *ChatPostgresRepository) Save(chat chat.Chat) (*chat.Chat, error) {
	stmt, err := rep.db.Prepare("INSERT INTO chats(id, first_user_id, second_user_id, created) VALUES($1, $2, $3, $4)")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to save chat, error: %w", err)
	}

	_, err = stmt.Exec(chat.Id.Value, chat.FirstUser.Value, chat.SecondUser.Value, chat.Created)
	if err != nil {
		return nil, fmt.Errorf("can't save chat with id = [%d], error: %w", chat.Id.Value, err)
	}
	return &chat, nil
}

func (rep *ChatPostgresRepository) Generate() (*chat.ChatId, error) {
	stmt, err := rep.db.Prepare("SELECT nextval('chat_id_seq')")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to get chat id, error: %w", err)
	}

	var id int64

	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("can't get new chat id, error: %w", err)
	}

	return &chat.ChatId{Value: id}, nil
}
