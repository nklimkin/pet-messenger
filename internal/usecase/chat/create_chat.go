package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type CreateChat interface {
	Execute(id chat.ChatId, firstUserId user.UserId, secondUserId user.UserId) (*chat.Chat, error)
}