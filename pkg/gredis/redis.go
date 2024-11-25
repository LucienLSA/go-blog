package gredis

import (
	"blog-service/global"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var RedisConn *redis.Pool

func RedisInit() error {
	// Redis配置
	RedisConn = &redis.Pool{
		MaxIdle:     global.RedisSetting.MaxIdle,
		MaxActive:   global.RedisSetting.MaxActive,
		IdleTimeout: global.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", global.RedisSetting.Host, global.RedisSetting.Port))
			if err != nil {
				return nil, err
			}
			if global.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", global.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err

				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Set(key string, data interface{}, expiration time.Duration) error {
	conn := RedisConn.Get()
	defer conn.Close()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, expiration)
	if err != nil {
		return err
	}
	return nil
}

// 判断key是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// Get 从redis中获取数据
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func Del(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", key+"*"))
	if err != nil {
		return err
	}
	for _, k := range keys {
		_, err = Del(k)
		if err != nil {
			return err
		}
	}
	return nil
}
