package user

import (
	"fmt"

	"ru.nklimkin/petmsngr/internal/domain/user"
)

type UserSignUpUseCase struct {
	userPersistence UserPersistence
	userIdGenerator UserIdGenerator
}

func New(userPersistence UserPersistence, userIdGenerator UserIdGenerator) *UserSignUpUseCase {
	return &UserSignUpUseCase{userPersistence: userPersistence, userIdGenerator: userIdGenerator}
}

func (uc *UserSignUpUseCase) Execute(login string) (*user.User, error) {
	newUserId, err := uc.userIdGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("can't handle request to sign up user, error: %w", err)
	}
	userToSave := user.User{Id: *newUserId, Login: login}
	return uc.userPersistence.Save(userToSave)
}
