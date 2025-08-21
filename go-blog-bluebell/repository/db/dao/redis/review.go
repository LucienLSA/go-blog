package redisCache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

const (
	KeyReviewTaskPrefix   = "review:task:"      // 审核任务前缀
	KeyReviewTaskQueue    = "review:task:queue" // 审核任务队列
	KeyReviewStatusPrefix = "review:status:"    // 审核状态前缀
)

// SaveReviewTask 保存审核任务到Redis
func SaveReviewTask(task *models.PostReviewTask) error {
	key := fmt.Sprintf("%s%d", KeyReviewTaskPrefix, task.PostID)

	// 将任务序列化为JSON
	jsonData, err := json.Marshal(task)
	if err != nil {
		zap.L().Error("marshal review task failed", zap.Error(err))
		return err
	}

	// 保存到Redis，设置24小时过期时间
	err = rdb.Set(rctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		zap.L().Error("save review task to redis failed", zap.Error(err))
		return err
	}

	// 添加到任务队列
	err = rdb.LPush(rctx, KeyReviewTaskQueue, jsonData).Err()
	if err != nil {
		zap.L().Error("add review task to queue failed", zap.Error(err))
		return err
	}

	zap.L().Info("save review task success",
		zap.Int64("post_id", task.PostID),
		zap.String("status", task.Status))

	return nil
}

// GetReviewTask 从Redis获取审核任务
func GetReviewTask(postID int64) (*models.PostReviewTask, error) {
	key := fmt.Sprintf("%s%d", KeyReviewTaskPrefix, postID)

	// 从Redis获取数据
	jsonData, err := rdb.Get(rctx, key).Result()
	if err != nil {
		zap.L().Error("get review task from redis failed", zap.Error(err))
		return nil, err
	}

	// 反序列化数据
	var task models.PostReviewTask
	err = json.Unmarshal([]byte(jsonData), &task)
	if err != nil {
		zap.L().Error("unmarshal review task failed", zap.Error(err))
		return nil, err
	}

	return &task, nil
}

// UpdateReviewTask 更新审核任务状态
func UpdateReviewTask(task *models.PostReviewTask) error {
	key := fmt.Sprintf("%s%d", KeyReviewTaskPrefix, task.PostID)

	// 将任务序列化为JSON
	jsonData, err := json.Marshal(task)
	if err != nil {
		zap.L().Error("marshal review task failed", zap.Error(err))
		return err
	}

	// 更新到Redis
	err = rdb.Set(rctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		zap.L().Error("update review task in redis failed", zap.Error(err))
		return err
	}

	zap.L().Info("update review task success",
		zap.Int64("post_id", task.PostID),
		zap.String("status", task.Status))

	return nil
}

// SaveReviewStatus 保存审核状态到Redis
func SaveReviewStatus(postID int64, status string, result *types.ReviewResult) error {
	key := fmt.Sprintf("%s%d", KeyReviewStatusPrefix, postID)

	statusData := map[string]interface{}{
		"status":     status,
		"result":     result,
		"updated_at": time.Now().Unix(),
	}

	// 序列化数据
	jsonData, err := json.Marshal(statusData)
	if err != nil {
		zap.L().Error("marshal review status failed", zap.Error(err))
		return err
	}

	// 保存到Redis，设置7天过期时间
	err = rdb.Set(rctx, key, jsonData, 7*24*time.Hour).Err()
	if err != nil {
		zap.L().Error("save review status to redis failed", zap.Error(err))
		return err
	}

	zap.L().Info("save review status success",
		zap.Int64("post_id", postID),
		zap.String("status", status))

	return nil
}

// GetReviewStatus 从Redis获取审核状态
func GetReviewStatus(postID int64) (string, *types.ReviewResult, error) {
	key := fmt.Sprintf("%s%d", KeyReviewStatusPrefix, postID)

	// 从Redis获取数据
	jsonData, err := rdb.Get(rctx, key).Result()
	if err != nil {
		zap.L().Error("get review status from redis failed", zap.Error(err))
		return "", nil, err
	}

	// 反序列化数据
	var statusData map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &statusData)
	if err != nil {
		zap.L().Error("unmarshal review status failed", zap.Error(err))
		return "", nil, err
	}

	status := statusData["status"].(string)
	var result *types.ReviewResult
	if resultData, ok := statusData["result"]; ok {
		resultBytes, _ := json.Marshal(resultData)
		json.Unmarshal(resultBytes, &result)
	}

	return status, result, nil
}

// GetNextReviewTask 从队列获取下一个审核任务
func GetNextReviewTask() (*models.PostReviewTask, error) {
	// 从队列右侧弹出任务
	result, err := rdb.BRPop(rctx, 5*time.Second, KeyReviewTaskQueue).Result()
	if err != nil {
		zap.L().Error("get next review task from queue failed", zap.Error(err))
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("no task in queue")
	}

	// 反序列化任务
	var task *models.PostReviewTask
	err = json.Unmarshal([]byte(result[1]), &task)
	if err != nil {
		zap.L().Error("unmarshal review task failed", zap.Error(err))
		return nil, err
	}

	return task, nil
}

// IsReviewTaskExists 检查审核任务是否存在
func IsReviewTaskExists(postID int64) bool {
	key := fmt.Sprintf("%s%d", KeyReviewTaskPrefix, postID)

	exists, err := rdb.Exists(rctx, key).Result()
	if err != nil {
		zap.L().Error("check review task exists failed", zap.Error(err))
		return false
	}

	return exists > 0
}
