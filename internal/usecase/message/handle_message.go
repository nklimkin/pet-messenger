package message

import "ru.nklimkin/petmsngr/internal/domain/message"

type HandleChatMessageUseCase struct {
	messageAccessor ChatMessageAccessor
	messagePersistence ChatMessagePersistence
}

func New(messageAccessor ChatMessageAccessor, messagePersistence ChatMessagePersistence) *HandleChatMessageUseCase {
	return &HandleChatMessageUseCase{
		messageAccessor: messageAccessor, 
		messagePersistence: messagePersistence}
}


func (us *HandleChatMessageUseCase) Handle(message message.NewMessage) string {	
	return ""
}