package models

type User struct {
	Age            uint8  `db:"age" json:"age"`
	UserID         int64  `db:"user_id" json:"user_id"`
	Username       string `db:"username" json:"user_name"`
	Password       string `db:"password"`
	Email          string `db:"email" json:"email"`
	PasswordDigest string
	Gender         string `db:"gender" json:"gender"`
	Token          string
	Avatar         string `db:"avatar" json:"avatar"`
}

type UserLogin struct {
	UserID   int64  `json:"user_id" db:"user_id"`   // 用户id
	Username string `json:"username" db:"username"` // 用户名
	Token    string `json:"token" `                 // access token
}
