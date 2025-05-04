package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserIdGenerator interface {
	Generate() (*user.UserId, error)
}