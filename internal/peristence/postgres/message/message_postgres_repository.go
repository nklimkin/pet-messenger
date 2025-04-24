package message

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/message"
)

type MessagePostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository() (*MessagePostgresRepository, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=messanger sslmode=disable password=postgres host=localhost")
	if err != nil {
		panic("Can't connect to db")
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS message (
		id BIGINT PRIMARY KEY,
		chat_id BIGINT PRIMARY KEY,
		status VARCHAR(255),
		payload TEXT,
		created TIMESTAMP);
	`)

	if err != nil {
		panic("Can't create message table")
	}

	_, err = stmt.Exec()
	if err != nil {
		panic("Can't create message table")
	}

	return &MessagePostgresRepository{db}, nil
}

func (rep *MessagePostgresRepository) GetByMessageId(id message.MessageId) message.ChatMessage {
	stmt, err := rep.db.Prepare("SELECT * FROM message WHERE id = $1")
	if err != nil {
		panic("Can't get message")
	}

	row := stmt.QueryRow(id.Value)

	return scanChatMessageRow(row)
}

func (rep *MessagePostgresRepository) GetByChatId(chatId chat.ChatId) []message.ChatMessage {
	stmt, err := rep.db.Prepare("SELECT * FROM message WHERE chat_id = $1 ORDER BY id")
	if err != nil {
		panic("Can't get message")
	}

	rows, err := stmt.Query(chatId.Value)
	if err != nil {
		panic("Can't get message")
	}

	messages := make([]message.ChatMessage, 1)

	for rows.Next() {
		rows.Scan()
		messages = append(messages, scanChatMessageRows(rows))
	}

	return messages
}

func scanChatMessageRow(row *sql.Row) message.ChatMessage {
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
		panic("Can't get message")
	}

	return buildChatMessage(persistedId, persistedChatId, persistedStatus, persistedPayload, persistedCreated)
}

func scanChatMessageRows(rows *sql.Rows) message.ChatMessage {
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
		panic("Can't get message")
	}

	return buildChatMessage(persistedId, persistedChatId, persistedStatus, persistedPayload, persistedCreated)
}

func buildChatMessage(
	persistedId int64,
	persistedChatId int64,
	persistedStatus string,
	persistedPayload string,
	persistedCreated time.Time,
) message.ChatMessage {
	messageId := message.MessageId{Value: persistedId}
	chatId := chat.ChatId{Value: persistedChatId}
	status := message.ParseMessageStatus(persistedStatus)
	payload := message.MessagePayload{Value: persistedPayload}
	return message.ChatMessage{Id: messageId,
		ChatId:  chatId,
		Status:  status,
		Payload: payload,
		Created: persistedCreated,
	}
}

func (rep *MessagePostgresRepository) Save(message message.ChatMessage) {
	stmt, err := rep.db.Prepare("INSERT INTO message(id, chat_id, status, payload, created) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		panic("Can't insert message")
	}

	_, err = stmt.Exec(message.Id.Value, message.ChatId.Value, message.Status.String(), message.Payload.Value, message.Created)
	if err != nil {
		panic("Can't insert message")
	}
}
