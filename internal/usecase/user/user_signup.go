package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserSignUp interface {
	Execute(id int64, login string) (*user.User, error)
}
