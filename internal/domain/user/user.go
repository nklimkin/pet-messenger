package user

type UserId struct {
	Value int64
}

type User struct {
	Id UserId
	Login string
}