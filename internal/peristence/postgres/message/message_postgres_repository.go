package message

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/config"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/message"

	_ "github.com/lib/pq"
)

const DATABASE_CONNECTION_TEMPLATE = "postgres://%s:%s@%s/%s?sslmode=disable"

type MessagePostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(datasource config.Datasource) (*MessagePostgresRepository, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			DATABASE_CONNECTION_TEMPLATE, datasource.Username, datasource.Password, datasource.Host, datasource.DatabaseName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS messages (
		id BIGINT PRIMARY KEY,
		chat_id BIGINT PRIMARY KEY,
		status VARCHAR(255),
		payload TEXT,
		created TIMESTAMP);
	`)

	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to create table - messages, error: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("can't create table - messages, error: %w", err)
	}

	return &MessagePostgresRepository{db}, nil
}

func (rep *MessagePostgresRepository) GetByMessageId(id message.MessageId) (*message.ChatMessage, error) {
	stmt, err := rep.db.Prepare("SELECT * FROM messages WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to get message with id = [%d], error: %w", id.Value, err)
	}

	row := stmt.QueryRow(id.Value)

	return scanChatMessageRow(row)
}

func (rep *MessagePostgresRepository) GetByChatId(chatId chat.ChatId) ([]*message.ChatMessage, error) {
	stmt, err := rep.db.Prepare("SELECT * FROM messages WHERE chat_id = $1 ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to get chat by id, error: %w", err)
	}

	rows, err := stmt.Query(chatId.Value)
	if err != nil {
		return nil, fmt.Errorf("can't get chat by id = [%d], error: %w", chatId.Value, err)
	}

	messages := make([]*message.ChatMessage, 1)

	for rows.Next() {
		rows.Scan()
		message, err := scanChatMessageRows(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func scanChatMessageRow(row *sql.Row) (*message.ChatMessage, error) {
	var persistedId int64
	var persistedChatId int64
	var persistedStatus string
	var persistedPayload string
	var persistedCreated time.Time

	err := row.Scan(&persistedId,
		&persistedChatId,
		&persistedStatus,
		&persistedPayload,
		&persistedCreated)
	if err != nil {
		return nil, fmt.Errorf("can't build message from result set: %w", err)
	}

	return buildChatMessage(persistedId, persistedChatId, persistedStatus, persistedPayload, persistedCreated), nil
}

func scanChatMessageRows(rows *sql.Rows) (*message.ChatMessage, error) {
	var persistedId int64
	var persistedChatId int64
	var persistedStatus string
	var persistedPayload string
	var persistedCreated time.Time

	err := rows.Scan(&persistedId,
		&persistedChatId,
		&persistedStatus,
		&persistedPayload,
		&persistedCreated)
	if err != nil {
		return nil, fmt.Errorf("can't build message from result set: %w", err)
	}

	return buildChatMessage(persistedId, persistedChatId, persistedStatus, persistedPayload, persistedCreated), nil
}

func buildChatMessage(
	persistedId int64,
	persistedChatId int64,
	persistedStatus string,
	persistedPayload string,
	persistedCreated time.Time,
) *message.ChatMessage {
	messageId := message.MessageId{Value: persistedId}
	chatId := chat.ChatId{Value: persistedChatId}
	status := message.ParseMessageStatus(persistedStatus)
	payload := message.MessagePayload{Value: persistedPayload}
	return &message.ChatMessage{Id: messageId,
		ChatId:  chatId,
		Status:  status,
		Payload: payload,
		Created: persistedCreated,
	}
}

func (rep *MessagePostgresRepository) Save(message message.ChatMessage) (*message.ChatMessage, error) {
	stmt, err := rep.db.Prepare("INSERT INTO messages(id, chat_id, status, payload, created) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, fmt.Errorf("can't build prepare statement to save message, error: %w", err)
	}

	_, err = stmt.Exec(message.Id.Value, message.ChatId.Value, message.Status.String(), message.Payload.Value, message.Created)
	if err != nil {
		return nil, fmt.Errorf("can't save message with id = [%d], error: %w", message.Id.Value, err)
	}
	return &message, nil
}
