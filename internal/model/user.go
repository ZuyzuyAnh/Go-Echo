package model

type User struct {
	ID       int64
	Email    string
	Name     string
	Password string
	Phone    string
}

type UserRole struct {
	UserID int64
	RoleID int64
}
