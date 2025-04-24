package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/user"
	"ru.nklimkin/petmsngr/internal/domain/chat"
)

type GetUserChats interface{
	Execute(userId user.UserId) []chat.Chat
}