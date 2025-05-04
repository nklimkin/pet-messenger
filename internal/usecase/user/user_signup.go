package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserSignUp interface {
	Execute(login string) (*user.User, error)
}
