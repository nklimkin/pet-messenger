package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserPersistence interface {
	Save(user user.User) (*user.User, error)
}