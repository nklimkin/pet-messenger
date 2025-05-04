package chat

import (
	"fmt"
	"sort"

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
	matchChats := make([]*chat.Chat, 0)
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
	rep.storage[chat.Id] = &chat
	return &chat, nil
}

func (rep *ChatInMemmoryRepository) Generate() (*chat.ChatId, error) {
	ids := make([]chat.ChatId, 0, len(rep.storage))

	for id := range rep.storage {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[j].Value < ids[i].Value
	})

	if len(ids) == 0 {
		return &chat.ChatId{Value: 1}, nil
	} else {
		return &ids[0], nil
	}
}
