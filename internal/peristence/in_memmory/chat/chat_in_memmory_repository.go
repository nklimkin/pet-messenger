package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type ChatInMemmoryRepository struct {
	storage map[chat.ChatId]chat.Chat
}

func New() *ChatInMemmoryRepository {
	return &ChatInMemmoryRepository{make(map[chat.ChatId]chat.Chat)}
}

func (rep *ChatInMemmoryRepository) GetById(id chat.ChatId) chat.Chat {
	return rep.storage[id]
}

func (rep *ChatInMemmoryRepository) GetByUserId(userId user.UserId) []chat.Chat {
	var matchChats []chat.Chat
	for _, item := range rep.storage {
		if item.FirstUser == userId || item.SecondUser == userId {
			matchChats = append(matchChats, item)
		}
	}
	return matchChats
}

func (rep *ChatInMemmoryRepository) Save(chat chat.Chat) {
	rep.storage[chat.Id] = chat
}
