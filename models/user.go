package models

type User struct {
	UserID         int64  `db:"user_id"`
	Username       string `db:"username"`
	Password       string `db:"password"`
	Email          string `db:"email"`
	Age            uint8  `db:"age"`
	PasswordDigest string
	Gender         string `db:"gender"`
}
