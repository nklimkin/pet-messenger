package user

import "ru.nklimkin/petmsngr/internal/domain/user"

type UserInMemmoryRepository struct {
	storage map[user.UserId]user.User
}

func New() *UserInMemmoryRepository {
	return &UserInMemmoryRepository{make(map[user.UserId]user.User)}
}

func (rep *UserInMemmoryRepository) GetById(id user.UserId) user.User {
	return rep.storage[id]
}

func (rep *UserInMemmoryRepository) Save(user user.User) {
	rep.storage[user.Id] = user
}