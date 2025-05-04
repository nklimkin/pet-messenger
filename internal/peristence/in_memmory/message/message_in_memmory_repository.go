package message

import (
	"fmt"
	"sort"

	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/message"
)

type MessageInMemmoryRepository struct {
	storage map[message.MessageId]*message.ChatMessage
}

func New() *MessageInMemmoryRepository {
	return &MessageInMemmoryRepository{make(map[message.MessageId]*message.ChatMessage)}
}

func (rep *MessageInMemmoryRepository) GetByMessageId(id message.MessageId) (*message.ChatMessage, error) {
	message := rep.storage[id]
	if message == nil {
		return nil, fmt.Errorf("there is no message with id = [%d]", id.Value)
	}
	return message, nil
}

func (rep *MessageInMemmoryRepository) GetByChatId(chatId chat.ChatId) ([]*message.ChatMessage, error) {
	var matchMessages []*message.ChatMessage
	for _, item := range rep.storage {
		if item.ChatId == chatId {
			matchMessages = append(matchMessages, item)
		}
	}
	if len(matchMessages) == 0 {
		return nil , fmt.Errorf("there is no messages for chat with id = [%d]", chatId.Value)
	}
	return matchMessages, nil
}

func (rep *MessageInMemmoryRepository) Save(message message.ChatMessage) (*message.ChatMessage, error) {
	rep.storage[message.Id] = &message
	return &message, nil
}

func (rep *MessageInMemmoryRepository) Generate() (*message.MessageId, error) {
	ids := make([]message.MessageId, 0, len(rep.storage))

	for id := range rep.storage {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[j].Value < ids[i].Value
	})

	if len(ids) == 0 {
		return &message.MessageId{Value: 1}, nil
	} else {
		return &ids[0], nil
	}
}
