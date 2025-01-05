package redis

import (
	"context"
	"fmt"

	"github.com/LucienLSA/go-blog/settings"
	"github.com/redis/go-redis/v9"
)

// var 声明全局的rdb变量
var rdb *redis.Client

// 初始化连接
func InitRedis(cfg *settings.RedisConfig) (err error) {
	// 不可在返回error时，新定义变量 会出现空指针问题，因为已经定义了全局变量
	rdb = redis.NewClient(&redis.Options{
		// Addr: fmt.Sprintf("%s:%d",
		// 	viper.GetString("redis.host"),
		// 	viper.GetInt("redis.port")),
		// Password: viper.GetString("redis.password"), // 密码
		// DB:       viper.GetInt("redis.db"),          // 数据库
		// PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
