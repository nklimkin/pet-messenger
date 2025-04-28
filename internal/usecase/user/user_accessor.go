package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserAccessor interface {
	GetById(id user.UserId) (*user.User, error)
}