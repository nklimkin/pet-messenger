package chat

import (
	"time"

	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type CreateChatUseCase struct {
	chatPerisistence ChatPersistence
}

func NewCreateChat(chatPersistence ChatPersistence) *CreateChatUseCase {
	return &CreateChatUseCase{chatPerisistence: chatPersistence}
}

func (uc *CreateChatUseCase) Execute(id chat.ChatId, fistUserId user.UserId, secondUserId user.UserId) (*chat.Chat, error) {
	newChat := chat.Chat{
		Id:         id,
		FirstUser:  fistUserId,
		SecondUser: secondUserId,
		Created:    time.Now(),
	}

	return uc.chatPerisistence.Save(newChat)
}
