package message

import "ru.nklimkin/petmsngr/internal/domain/message"

type HandleChatMessageUseCase struct {
	chatAccessor ChatMessageAccessor
	chatPersistence ChatMessagePersistence
}

func New(chatAccessor ChatMessageAccessor, chatPersistence ChatMessagePersistence) *HandleChatMessageUseCase {
	return &HandleChatMessageUseCase{chatAccessor: chatAccessor, chatPersistence: chatPersistence}
}


func (us *HandleChatMessageUseCase) Handle(message message.NewMessage) string {	
	return ""
}