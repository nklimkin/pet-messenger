package message

import (
	"ru.nklimkin/petmsngr/internal/domain/message"
)

type ChatMessagePersistence interface {
	Save(message message.ChatMessage) (*message.ChatMessage, error)
}
