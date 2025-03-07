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
	ErrNotExistToken = errors.New("不存在的key")
)

// Tokens存入redis
func StorgeUserIdToken(token, username string) (err error) {
	// 获取token的存活实践
	duration := time.Duration(settings.Conf.AppConfig.JwtExpireTime * time.Hour)
	key := GetRedisKey(KeyUserIDTokenSetPrefix)
	// 存入redis
	if err = rdb.Set(rctx, key+username, token, duration).Err(); err != nil {
		zap.L().Error("Insert username, token into redis failed, err:%v", zap.Error(err))
		// zap.L().Debug("Insert username, token into redis failed, err:%v", zap.Error(err))
		fmt.Printf("Insert username, token into redis failed, err:%v", zap.Error(err))
		return
	}
	return err
}

// 从redis取token
func GetJwtToken(username string) (token string, err error) {
	key := GetRedisKey(KeyUserIDTokenSetPrefix)
	token, err = rdb.Get(rctx, key+username).Result()
	if err == redis.Nil {
		return "", ErrNotExistToken
	}
	if err != nil {
		zap.L().Error("Search token from redis failed, err:%v", zap.Error(err))
		fmt.Printf("Search token from redis failed, err:%v", zap.Error(err))
		return
	}
	return
}
