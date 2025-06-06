package message

import (
	"ru.nklimkin/petmsngr/internal/domain/message"
	"ru.nklimkin/petmsngr/internal/domain/chat"
)

type ChatMessageAccessor interface {
	GetByMessageId(id message.MessageId) (*message.ChatMessage, error)
	GetByChatId(chatId chat.ChatId) ([]*message.ChatMessage, error)
}
