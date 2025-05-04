package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type CreateChat interface {
	Execute(firstUserId user.UserId, secondUserId user.UserId) (*chat.Chat, error)
}