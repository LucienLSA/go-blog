package redisCache

import (
	"errors"
	"fmt"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"go.uber.org/zap"
)

var (
	ErrExistCode = errors.New("邮箱验证码已发送过")
	ErrEmailCode = errors.New("邮箱验证码不正确")
)

// 邮箱和验证码存入redis
func StorgeEmailCode(email, code string) (err error) {
	duration := time.Duration(settings.Conf.EmailConfig.EmailCodeSaveTime * time.Minute)
	key := GetRedisKey(KeyEmailSetPrefix)
	err = rdb.Set(rctx, key+email, code, duration).Err()
	if err != nil {
		zap.L().Error("Insert email code into redis failed, err:%v\n", zap.Error(err))
		fmt.Printf("Insert email code into redis failed, err:%v\n", err)
		return err
	}
	return
}

// 从redis查看邮箱和验证码是否存在
func EmailCodeExists(email string) (err error) {
	key := GetRedisKey(KeySendEmailSetPrefix)
	if rdb.Exists(rctx, key+email).Val() > 0 {
		zap.L().Error("Get email code from redis")
		fmt.Printf("Get email code from redis")
		return ErrExistCode
	}
	return
}

// 发送邮箱和验证码存入redis，区别于StorgeEmailCode是为了验证发送验证码的时效
func SendEmailCode(email, code string) (err error) {
	duration := time.Duration(settings.Conf.EmailConfig.EmailCodeExpireTime * time.Second)
	key := GetRedisKey(KeySendEmailSetPrefix)
	err = rdb.Set(rctx, key+email, code, duration).Err()
	if err != nil {
		zap.L().Error("Insert send-email code into redis failed, err:%v\n", zap.Error(err))
		fmt.Printf("Insert send-email code into redis failed, err:%v\n", err)
		return err
	}
	return
}

// // 检验redis中邮箱验证码的正确性
// func CheckEmailCode(code int32) (err error) {
// 	key1 := GetRedisKey(KeyEmailSetPrefix)
// 	key2 := GetRedisKey(KeySendEmailSetPrefix)
// 	codeStroage := rdb.Get(rctx, key1+email).Val()
// 	codeSend := rdb.Get(rctx, key2+email).Val()
// 	if codeStroage != codeSend {
// 		zap.L().Error("check email code failed")
// 		fmt.Printf("check email code failed")
// 		return ErrEmailCode
// 	}
// 	return
// }
