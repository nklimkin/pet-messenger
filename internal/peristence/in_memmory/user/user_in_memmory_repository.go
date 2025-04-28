package user

import (
	"fmt"

	"ru.nklimkin/petmsngr/internal/domain/user"
)

type UserInMemmoryRepository struct {
	storage map[user.UserId]*user.User
}

func New() *UserInMemmoryRepository {
	return &UserInMemmoryRepository{make(map[user.UserId]*user.User)}
}

func (rep *UserInMemmoryRepository) GetById(id user.UserId) (*user.User, error) {
	user := rep.storage[id]
	if user == nil {
		return nil, fmt.Errorf("there is no user with id [%d]", id.Value)
	}
	return rep.storage[id], nil
}

func (rep *UserInMemmoryRepository) Save(user user.User) (*user.User, error) {
	rep.storage[user.Id] = &user
	return &user, nil
}