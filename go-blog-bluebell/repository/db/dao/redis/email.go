package redisCache

import (
	"errors"
	"fmt"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrExistCode   = errors.New("邮箱验证码已发送过")
	ErrEmailCode   = errors.New("邮箱验证码不正确")
	ErrCodeExpired = errors.New("验证码已过期")
)

// 邮箱和验证码存入redis
func StorgeEmailCode(email string, code string, opType int) error {
	duration := time.Duration(settings.Conf.EmailConfig.EmailCodeSaveTime * time.Minute)
	key := fmt.Sprintf("email_code:%s:%d", email, opType)
	err := rdb.Set(rctx, key, code, duration).Err()
	if err != nil {
		zap.L().Error("Insert email code into redis failed", zap.Error(err))
		return err
	}
	return nil
}

// 从redis查看邮箱和验证码是否存在
func EmailCodeExists(email string, opType int) error {
	key := fmt.Sprintf("email_code:%s:%d", email, opType)
	if rdb.Exists(rctx, key).Val() > 0 {
		return ErrExistCode
	}
	return nil
}

// 发送邮箱和验证码存入redis，区别于StorgeEmailCode是为了验证发送验证码的时效
func SendEmailCode(email string, code string, opType int) error {
	duration := time.Duration(settings.Conf.EmailConfig.EmailCodeExpireTime * time.Second)
	key := fmt.Sprintf("email_code:%s:%d", email, opType)
	err := rdb.Set(rctx, key, code, duration).Err()
	if err != nil {
		zap.L().Error("Insert send-email code into redis failed", zap.Error(err))
		return err
	}
	return nil
}

// 检验redis中邮箱验证码的正确性
func CheckEmailCode(email string, code string, opType int) error {
	key := fmt.Sprintf("email_code:%s:%d", email, opType)
	codeStored, err := rdb.Get(rctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCodeExpired
		}
		zap.L().Error("CheckEmailCode failed", zap.Error(err), zap.String("email", email), zap.Int("opType", opType))
		return errors.New("验证码验证失败，请重试")
	}
	if codeStored != code {
		return ErrEmailCode
	}
	// 验证通过后删除验证码
	rdb.Del(rctx, key)
	return nil
}
