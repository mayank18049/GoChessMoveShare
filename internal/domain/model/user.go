package model

func (u *User) Equals(userID string) bool {
	return u.id == userID
}

func NewUser(id string) (*User, UsersStatus) {
	return &User{id: id}, USER_OK
}

func (u Users) ContainsUserID(userID string) bool {
	for _, user := range u {
		if user.Equals(userID) {
			return true
		}
	}
	return false
}

func (u *Users) AddUser(user *User) UsersStatus {
	if (*u).ContainsUserID(user.id) {
		return USER_EXISTS
	}
	*u = append((*u), user)
	return USER_OK
}

func (u User) GetID() string {
	return u.id
}
