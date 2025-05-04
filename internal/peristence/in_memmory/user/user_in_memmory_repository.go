package user

import (
	"fmt"
	"sort"

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

func (rep *UserInMemmoryRepository) Generate() (*user.UserId, error) {
	ids := make([]user.UserId, 0, len(rep.storage))

	for id := range rep.storage {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[j].Value < ids[i].Value
	})

	if len(ids) == 0 {
		return &user.UserId{Value: 1}, nil
	} else {
		return &ids[0], nil
	}
}

func (rep *UserInMemmoryRepository) Exists(id user.UserId) (bool, error) {
	return rep.storage[id] != nil, nil
}
