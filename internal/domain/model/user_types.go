package model

type UsersStatus int

const (
	USER_OK UsersStatus = -iota
	USER_EXISTS
	USER_NOT_FOUND
	OUT_OF_BOUNDS
	INVALID_USERID
)

type User struct {
	id string
}

type Users []*User
