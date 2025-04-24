package message

import (
	"ru.nklimkin/petmsngr/internal/domain/message"
	"ru.nklimkin/petmsngr/internal/domain/chat"
)

type MessageInMemmoryRepository struct {
	storage map[message.MessageId]message.ChatMessage
}

func New() *MessageInMemmoryRepository {
	return &MessageInMemmoryRepository{make(map[message.MessageId]message.ChatMessage)}
}

func (rep *MessageInMemmoryRepository) GetByMessageId(id message.MessageId) message.ChatMessage {
	return rep.storage[id]
}

func (rep *MessageInMemmoryRepository) GetByChatId(chatId chat.ChatId) []message.ChatMessage {
	var matchMessages []message.ChatMessage
	for _, item := range rep.storage {
		if item.ChatId == chatId {
			matchMessages = append(matchMessages, item)
		}
	}
	return matchMessages
}

func (rep *MessageInMemmoryRepository) Save(message message.ChatMessage) {
	rep.storage[message.Id] = message
}
