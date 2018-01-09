package domain

type UserID uint64

// User is the representation of a single user.
type User struct {
	ID       UserID
	Nickname string
	RealName string
}

// UserWithPassword is a user with embedded password field.
// Typically used for UserCreate request only.
type UserWithPassword struct {
	User
	Password string
}
