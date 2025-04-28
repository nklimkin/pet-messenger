package user

import (
	"ru.nklimkin/petmsngr/internal/domain/user"
)

type UserSignUpUseCase struct {
	userPersistence UserPersistence
}

func New(userPersistence UserPersistence) *UserSignUpUseCase {
	return &UserSignUpUseCase{userPersistence: userPersistence}
}

func (uc *UserSignUpUseCase) Execute(id int64, login string) (*user.User, error) {
	userId := user.UserId{Value: id}
	userToSave := user.User{Id: userId, Login: login}
	return uc.userPersistence.Save(userToSave)
}
