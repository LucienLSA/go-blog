package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorEmailNotExit    = errors.New("邮箱不存在")
	ErrorEmailExist      = errors.New("邮箱已存在")
	ErrorInvalidPassword = errors.New("密码错误")

	ErrorInvalidID = errors.New("无效的ID")
	// Error
)
