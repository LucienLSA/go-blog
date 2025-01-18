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

type UserLogin struct {
	UserID   int64  `json:"user_id" db:"user_id"`    // 用户id
	Username string `json:"user_name" db:"username"` // 用户名
	Token    string `json:"token" `                  // access token
}
