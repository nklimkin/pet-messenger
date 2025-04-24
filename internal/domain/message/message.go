package message

import (
	"time"
	"ru.nklimkin/petmsngr/internal/domain/chat"
)

type MessageId struct {
	Value int64
}

type MessagePayload struct {
	Value string
}

type MessageStatus int 

const (
	NEW MessageStatus = iota
	SENT
	READ
	UNKNOWN
)

var messageNames = map[MessageStatus]string{
	NEW: "NEW",
	SENT: "SENT",
	READ: "READ",
	UNKNOWN: "UNKNOWN",
}

func (s MessageStatus) String() string {
	return messageNames[s]
}

func ParseMessageStatus(source string) MessageStatus {
	for status, value := range messageNames {
		if value == source {
			return status
		}
	}
	return UNKNOWN
}

type ChatMessage struct {
	Id MessageId
	ChatId chat.ChatId
	Payload MessagePayload
	Status MessageStatus
	Created time.Time
}

type NewMessage struct {
	ChatId  chat.ChatId
	Payload MessagePayload
}
