package user

type UserSignUp interface {
	Execute(id int64, login string)
}
