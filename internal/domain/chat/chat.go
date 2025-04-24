package chat

import (
	"time"
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type ChatId struct {
	Value int64
}

type Chat struct {
	Id ChatId
	FirstUser user.UserId
	SecondUser user.UserId
	Created time.Time
}