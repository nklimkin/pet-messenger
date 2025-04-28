package chat

import "ru.nklimkin/petmsngr/internal/domain/chat"

type ChatPersistence interface {
	Save(chat chat.Chat) (*chat.Chat, error)
}