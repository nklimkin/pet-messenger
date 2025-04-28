package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type ChatAccessor interface {
	GetById(id chat.ChatId) (*chat.Chat, error) 
	GetByUserId(userId user.UserId) ([]*chat.Chat, error)
}