package models

type User struct {
	Age            uint8  `db:"age"`
	UserID         int64  `db:"user_id"`
	Username       string `db:"username"`
	Password       string `db:"password"`
	Email          string `db:"email"`
	PasswordDigest string
	Gender         string `db:"gender"`
	Token          string
}
