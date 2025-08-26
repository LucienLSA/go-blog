package redisCache

import (
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	KeyOAuthStatePrefix = "oauth:state:"
)

// SaveOAuthState stores a short-lived OAuth state to prevent CSRF
// 生成授权请求时，将随机生成的state值存入 Redis，设置短期过期时间（ttl，通常 5-10 分钟）。
func SaveOAuthState(state string, ttl time.Duration) error {
	key := KeyOAuthStatePrefix + state
	if err := rdb.Set(rctx, key, "1", ttl).Err(); err != nil {
		zap.L().Error("save oauth state failed", zap.Error(err))
		return err
	}
	return nil
}

// ConsumeOAuthState validates and deletes a state value. Returns true if valid.
// 第三方服务（如 GitHub）回调时，验证传入的state是否有效（存在于 Redis 且未过期），
// 验证通过后立即删除该state（防止重复使用）
func ConsumeOAuthState(state string) (bool, error) {
	key := KeyOAuthStatePrefix + state
	// Use Lua script for atomic GET + DEL to support older Redis without GETDEL
	script := `local v = redis.call('GET', KEYS[1]); if v then redis.call('DEL', KEYS[1]); end; return v`
	res, err := rdb.Eval(rctx, script, []string{key}).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		zap.L().Error("consume oauth state failed", zap.Error(err))
		return false, err
	}
	if res == nil {
		return false, nil
	}
	if s, ok := res.(string); ok {
		return s == "1", nil
	}
	return false, nil
}
