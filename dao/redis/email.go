package redisCache

import (
	"errors"
)

var (
	ErrNotExistCode = errors.New("邮箱验证码不存在")
	ErrEmailCode    = errors.New("邮箱验证码不正确")
)

// 邮箱验证码存入redis
func StorgeEmailCode() (err error) {

	return err
}

// 从redis取邮箱验证码
func GetEmailCode(email string) (token string, err error) {
	rdb.Get(rctx)
	return
}
