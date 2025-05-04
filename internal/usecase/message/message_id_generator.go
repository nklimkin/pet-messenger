package message

import "ru.nklimkin/petmsngr/internal/domain/message"

type MessageIdGenerator interface {
	Generate() (*message.MessageId, error)
}