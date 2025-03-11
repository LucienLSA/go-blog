package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID         int64  `gorm:"user_id"`
	Age            uint8  `gorm:"age"`
	UserName       string `gorm:"unique"`
	Password       string `gorm:"password"`
	Email          string `gorm:"email"`
	PasswordDigest string `gorm:"password_digest"`
	Gender         string `gorm:"gender"`
	Token          string `gorm:"token"`
	Avatar         string `gorm:"avatar"`
}

// 实现TableName方法可以修改表名
func (u *User) TableName() string {
	return "user"
}

// 用户登录请求
type UserLogin struct {
	UserID   int64  `json:"user_id" gorm:"user_id"`     // 用户id
	Username string `json:"user_name" gorm:"user_name"` // 用户名
	Token    string `json:"token" `                     // access token
}

// // 密码加密 旧版
// func encryptPassword(oPassword string) string {
// 	h := md5.New()
// 	h.Write([]byte(secret))
// 	return hex.EncodeToString(h.Sum([]byte(oPassword)))
// }

const (
	PasswordCost = 12 //密码加密难度
)

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	// 将原密码进行加密
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
