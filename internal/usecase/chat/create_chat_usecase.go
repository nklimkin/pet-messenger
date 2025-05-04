package chat

import (
	"fmt"
	"time"

	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
	user_usecase "ru.nklimkin/petmsngr/internal/usecase/user"
)

type CreateChatUseCase struct {
	chatPerisistence ChatPersistence
	chatIdGenerator  ChatIdGenerator
	userAccessor     user_usecase.UserAccessor
}

func NewCreateChat(
	chatPersistence ChatPersistence,
	chatIdGenerator ChatIdGenerator,
	userAccessor user_usecase.UserAccessor,
) *CreateChatUseCase {
	return &CreateChatUseCase{
		chatPerisistence: chatPersistence,
		chatIdGenerator:  chatIdGenerator,
		userAccessor:     userAccessor,
	}
}

func (uc *CreateChatUseCase) Execute(fistUserId user.UserId, secondUserId user.UserId) (*chat.Chat, error) {
	err := uc.checkUserExists(fistUserId)
	if err != nil {
		return nil, err
	}
	err = uc.checkUserExists(secondUserId)
	if err != nil {
		return nil, err
	}
	newChatId, err := uc.chatIdGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("can't handle request to create new chat, error: %w", err)
	}
	newChat := chat.Chat{
		Id:         *newChatId,
		FirstUser:  fistUserId,
		SecondUser: secondUserId,
		Created:    time.Now(),
	}

	return uc.chatPerisistence.Save(newChat)
}

func (uc *CreateChatUseCase) checkUserExists(userId user.UserId) error {
	isExixsts, err := uc.userAccessor.Exists(userId)
	if err != nil {
		return fmt.Errorf("error while check user existence, error: %w", err)
	}

	if isExixsts == false {
		return fmt.Errorf("there is no user with id = %d", userId.Value)
	}

	return nil
}
