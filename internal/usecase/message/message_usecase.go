package message

import (
	"ru.nklimkin/petmsngr/internal/domain/message"
)

type HandleChatMessage interface {
	Handle(message message.NewMessage) string
}
