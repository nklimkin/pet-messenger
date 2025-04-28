package chat

import (
	"fmt"

	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type ChatInMemmoryRepository struct {
	storage map[chat.ChatId]*chat.Chat
}

func New() *ChatInMemmoryRepository {
	return &ChatInMemmoryRepository{make(map[chat.ChatId]*chat.Chat)}
}

func (rep *ChatInMemmoryRepository) GetById(id chat.ChatId) (*chat.Chat, error) {
	chat := rep.storage[id]
	if chat == nil {
		return nil, fmt.Errorf("there is no chat with id = [%d]", id.Value)
	}
	return chat, nil
}

func (rep *ChatInMemmoryRepository) GetByUserId(userId user.UserId) ([]*chat.Chat, error) {
	matchChats []*chat.Chat := make(*[]chat.Chat)
	for _, item := range rep.storage {
		if item.FirstUser == userId || item.SecondUser == userId {
			matchChats = append(matchChats, item)
		}
	}
	if len(matchChats) == 0 {
		return nil, fmt.Errorf("there is no chats for user id = [%d]", userId.Value)
	}
	return matchChats, nil
}

func (rep *ChatInMemmoryRepository) Save(chat chat.Chat) (*chat.Chat, error) {
	rep.storage[chat.Id] = chat
	return &chat
}
