package chat

import (
	"time"

	"github.com/jmoiron/sqlx"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type ChatPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository() (*ChatPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=messanger sslmode=disable password=postgres host=localhost")
	if err != nil {
		panic("Can't connect to db")
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS chats (
		id BIGINT PRIMARY KEY,
		first_user_id BIGINT,
		second_user_id BIGINT,
		created TIMESTAMP);
	`)

	if err != nil {
		panic("Can't create chats table")
	}

	_, err = stmt.Exec()
	if err != nil {
		panic("Can't create chats table")
	}

	return &ChatPostgresRepository{db}, nil
}

func (rep *ChatPostgresRepository) GetById(id chat.ChatId) chat.Chat {
	stmt, err := rep.db.Prepare("SELECT * FROM chats WHERE id = $1")
	if err != nil {
		panic("Can't get chat")
	}

	var persistedId int64
	var persistedFirstUserId int64
	var persistedSecondUserId int64
	var persistedCreated time.Time

	err = stmt.QueryRow(id.Value).Scan(&persistedId, &persistedFirstUserId, &persistedSecondUserId, &persistedCreated)
	if err != nil {
		panic("Can't get chat")
	}

	return buildChat(persistedId, persistedFirstUserId, persistedSecondUserId, persistedCreated)
}

func (rep *ChatPostgresRepository) GetByUserId(userId user.UserId) []chat.Chat {
	stmt, err := rep.db.Prepare("SELECT * FROM chats WHERE first_user_id = $1 OR second_user_id = $1")
	if err != nil {
		panic("Can't get chats")
	}

	rows, err := stmt.Query(userId.Value)
	if err != nil {
		panic("Can't get chats")
	}

	chats := make([]chat.Chat, 1)

	for rows.Next() {
		var persistedId int64
		var persistedFirstUserId int64
		var persistedSecondUserId int64
		var persistedCreated time.Time

		err = rows.Scan(&persistedId, &persistedFirstUserId, &persistedSecondUserId, &persistedCreated)
		if err != nil {
			panic("Can't get chats")
		}

		chats = append(chats, buildChat(persistedId, persistedFirstUserId, persistedSecondUserId, persistedCreated))
	}

	return chats
}

func buildChat(
	persistedId int64, 
	persistedFirstUserId int64, 
	persistedSecondUserId int64, 
	persistedCreated time.Time,
	) chat.Chat {
	chatId := chat.ChatId{Value: persistedId}
	firstUserId := user.UserId{Value: persistedFirstUserId}
	secondUserId := user.UserId{Value: persistedSecondUserId}

	return chat.Chat{Id: chatId, FirstUser: firstUserId, SecondUser: secondUserId, Created: persistedCreated}
}

func (rep *ChatPostgresRepository) Save(chat chat.Chat) {
	stmt, err := rep.db.Prepare("INSERT INTO chats(id, first_user_id, second_user_id, created) VALUES($1, $2, $3, $4)")
	if err != nil {
		panic("Can't save chat")
	}

	_, err = stmt.Exec(chat.Id.Value, chat.FirstUser.Value, chat.SecondUser.Value, chat.Created)
	if err != nil {
		panic("Can't save chat")
	}

}