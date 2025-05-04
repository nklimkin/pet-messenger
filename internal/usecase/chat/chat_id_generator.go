package chat

import "ru.nklimkin/petmsngr/internal/domain/chat"

type ChatIdGenerator interface {
	Generate() (*chat.ChatId, error)
}