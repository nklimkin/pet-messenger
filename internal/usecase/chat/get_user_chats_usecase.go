package chat

import (
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type GetUserChatsUseCase struct {
	chatAccessor ChatAccessor
}

func NewGetUserChats(chatAccessor ChatAccessor) *GetUserChatsUseCase {
	return &GetUserChatsUseCase{chatAccessor: chatAccessor}
}

func (uc *GetUserChatsUseCase) Execute(userId user.UserId) ([]*chat.Chat, error) {
	return uc.chatAccessor.GetByUserId(userId)
}