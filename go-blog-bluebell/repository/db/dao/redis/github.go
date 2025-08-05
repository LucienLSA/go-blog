package redisCache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

const (
	KeyGitHubTrendingPrefix = "github:trending:"     // GitHub热点数据前缀
	KeyGitHubTrendingHash   = "github:trending:hash" // GitHub热点数据哈希表
)

// SaveGitHubTrendingData 保存GitHub热点数据到Redis
func SaveGitHubTrendingData(language, since string, data []types.GitHubTrendingData) error {
	key := fmt.Sprintf("%s%s:%s", KeyGitHubTrendingPrefix, language, since)

	// 将数据序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("marshal github trending data failed", zap.Error(err))
		return err
	}

	// 保存到Redis，设置1小时过期时间
	err = rdb.Set(rctx, key, jsonData, time.Hour).Err()
	if err != nil {
		zap.L().Error("save github trending data to redis failed", zap.Error(err))
		return err
	}

	// 同时保存到哈希表中，用于记录所有语言和时间范围的数据
	hashKey := fmt.Sprintf("%s:%s", language, since)
	err = rdb.HSet(rctx, KeyGitHubTrendingHash, hashKey, jsonData).Err()
	if err != nil {
		zap.L().Error("save github trending data to hash failed", zap.Error(err))
		return err
	}

	// 设置哈希表过期时间为1小时
	rdb.Expire(rctx, KeyGitHubTrendingHash, time.Hour)

	zap.L().Info("save github trending data success",
		zap.String("language", language),
		zap.String("since", since),
		zap.Int("count", len(data)))

	return nil
}

// GetGitHubTrendingData 从Redis获取GitHub热点数据
func GetGitHubTrendingData(language, since string) ([]types.GitHubTrendingData, error) {
	key := fmt.Sprintf("%s%s:%s", KeyGitHubTrendingPrefix, language, since)

	// 从Redis获取数据
	jsonData, err := rdb.Get(rctx, key).Result()
	if err != nil {
		zap.L().Error("get github trending data from redis failed", zap.Error(err))
		return nil, err
	}

	// 反序列化数据
	var data []types.GitHubTrendingData
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		zap.L().Error("unmarshal github trending data failed", zap.Error(err))
		return nil, err
	}

	return data, nil
}

// GetAllGitHubTrendingData 获取所有GitHub热点数据
func GetAllGitHubTrendingData() (map[string][]types.GitHubTrendingData, error) {
	// 从哈希表获取所有数据
	hashData, err := rdb.HGetAll(rctx, KeyGitHubTrendingHash).Result()
	if err != nil {
		zap.L().Error("get all github trending data from hash failed", zap.Error(err))
		return nil, err
	}

	result := make(map[string][]types.GitHubTrendingData)
	for key, jsonData := range hashData {
		var data []types.GitHubTrendingData
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			zap.L().Error("unmarshal github trending data failed", zap.Error(err))
			continue
		}
		result[key] = data
	}

	return result, nil
}

// IsGitHubTrendingDataExpired 检查GitHub热点数据是否过期
func IsGitHubTrendingDataExpired(language, since string) bool {
	key := fmt.Sprintf("%s%s:%s", KeyGitHubTrendingPrefix, language, since)

	// 检查key是否存在
	exists, err := rdb.Exists(rctx, key).Result()
	if err != nil {
		zap.L().Error("check github trending data expired failed", zap.Error(err))
		return true
	}

	return exists == 0
}
