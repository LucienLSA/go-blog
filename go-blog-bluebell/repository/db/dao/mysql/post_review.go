package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostReviewDao struct {
	*gorm.DB
}

func NewPostReviewDao(ctx context.Context) *PostReviewDao {
	return &PostReviewDao{NewDBClient(ctx)}
}

func NewPostReviewDaoByDB(db *gorm.DB) *PostReviewDao {
	return &PostReviewDao{db}
}

// CreateReviewTask 创建审核任务
func (dao *PostReviewDao) CreateReviewTask(task *models.PostReviewTask) error {
	return dao.DB.Model(&models.PostReviewTask{}).Create(task).Error
}

// GetReviewTaskByPostID 根据帖子ID获取审核任务
func (dao *PostReviewDao) GetReviewTaskByPostID(postID int64) (*models.PostReviewTask, error) {
	var task models.PostReviewTask
	err := dao.DB.Model(&models.PostReviewTask{}).Where("post_id = ?", postID).First(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateReviewTask 更新审核任务
func (dao *PostReviewDao) UpdateReviewTask(task *models.PostReviewTask) error {
	return dao.DB.Model(&models.PostReviewTask{}).Where("id = ?", task.ID).Updates(task).Error
}

// GetPendingReviewTasks 获取待处理的审核任务
func (dao *PostReviewDao) GetPendingReviewTasks(limit int) ([]*models.PostReviewTask, error) {
	var tasks []*models.PostReviewTask
	err := dao.DB.Model(&models.PostReviewTask{}).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// UpdateReviewTaskStatus 更新审核任务状态
func (dao *PostReviewDao) UpdateReviewTaskStatus(taskID int64, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	return dao.DB.Model(&models.PostReviewTask{}).Where("id = ?", taskID).Updates(updates).Error
}

// CreateReviewLog 创建审核日志
func (dao *PostReviewDao) CreateReviewLog(log *models.PostReviewLog) error {
	return dao.DB.Model(&models.PostReviewLog{}).Create(log).Error
}

// GetReviewLogsByPostID 根据帖子ID获取审核日志
func (dao *PostReviewDao) GetReviewLogsByPostID(postID int64) ([]*models.PostReviewLog, error) {
	var logs []*models.PostReviewLog
	err := dao.DB.Model(&models.PostReviewLog{}).
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// UpdatePostReviewStatus 更新帖子审核状态
func (dao *PostReviewDao) UpdatePostReviewStatus(postID int64, reviewResult *types.ReviewResult) error {
	// 将审核结果序列化为JSON
	resultJSON, err := json.Marshal(reviewResult)
	if err != nil {
		zap.L().Error("marshal review result failed", zap.Error(err))
		return err
	}

	// 将标签序列化为JSON
	tagsJSON, err := json.Marshal(reviewResult.Tags)
	if err != nil {
		zap.L().Error("marshal review tags failed", zap.Error(err))
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	updates := map[string]interface{}{
		"review_status":      reviewResult.Status,
		"review_result":      string(resultJSON),
		"reviewed_at":        &now,
		"review_score":       reviewResult.Score,
		"review_reason":      reviewResult.Reason,
		"review_suggestions": reviewResult.Suggestions,
		"review_tags":        string(tagsJSON),
	}

	return dao.DB.Model(&models.Post{}).Where("post_id = ?", postID).Updates(updates).Error
}

// GetPostReviewStatus 获取帖子审核状态
func (dao *PostReviewDao) GetPostReviewStatus(postID int64) (*models.Post, error) {
	var post models.Post
	err := dao.DB.Model(&models.Post{}).Where("post_id = ?", postID).First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostsByReviewStatus 根据审核状态获取帖子列表
func (dao *PostReviewDao) GetPostsByReviewStatus(status string, pageNum, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (pageNum - 1) * pageSize

	err := dao.DB.Model(&models.Post{}).
		Where("review_status = ?", status).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&posts).Error

	if err == sql.ErrNoRows {
		zap.L().Warn("no posts found with review status", zap.String("status", status))
		return nil, nil
	}

	return posts, err
}

// GetReviewStatistics 获取审核统计信息
func (dao *PostReviewDao) GetReviewStatistics() (map[string]int64, error) {
	var stats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	err := dao.DB.Model(&models.Post{}).
		Select("review_status as status, count(*) as count").
		Group("review_status").
		Find(&stats).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, stat := range stats {
		result[stat.Status] = stat.Count
	}

	return result, nil
}

// GetFailedReviewTasks 获取失败的审核任务
func (dao *PostReviewDao) GetFailedReviewTasks(limit int) ([]*models.PostReviewTask, error) {
	var tasks []*models.PostReviewTask
	err := dao.DB.Model(&models.PostReviewTask{}).
		Where("status = ?", "failed").
		Order("created_at DESC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// RetryFailedReviewTask 重试失败的审核任务
func (dao *PostReviewDao) RetryFailedReviewTask(taskID int64) error {
	updates := map[string]interface{}{
		"status":        "pending",
		"updated_at":    time.Now(),
		"completed_at":  nil,
		"error_message": "",
	}

	return dao.DB.Model(&models.PostReviewTask{}).Where("id = ?", taskID).Updates(updates).Error
}
