package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// 单例模式
var reviewSrvIns *ReviewSrv
var reviewSrvOnce sync.Once

type ReviewSrv struct {
}

// 单例实例
func GetReviewSrv() *ReviewSrv {
	reviewSrvOnce.Do(func() {
		reviewSrvIns = &ReviewSrv{}
	})
	return reviewSrvIns
}

// 重置单例 便于测试
func ResetReviewSrv() {
	reviewSrvOnce = sync.Once{}
	reviewSrvIns = nil
}

// SubmitPostForReview 提交帖子进行审核
func (s *ReviewSrv) SubmitPostForReview(ctx context.Context, req *types.PostReviewRequest) (*types.PostReviewResponse, error) {
	// 获取当前用户信息
	user, err := ctl.GetUserInfo(ctx)
	if err != nil {
		zap.L().Error("get user info failed", zap.Error(err))
		return nil, err
	}

	// 生成帖子ID
	postID := int64(snowflake.GenID())

	// 创建帖子（待审核状态）
	postDao := mysql.NewPostDao(ctx)
	post := &models.Post{
		Status:       1, // 活跃状态
		PostID:       postID,
		CommunityID:  req.CommunityID,
		AuthorID:     user.UserId,
		Title:        req.Title,
		Content:      req.Content,
		ReviewStatus: "pending",
		ReviewScore:  0.0,
	}

	err = postDao.CreatePostWithReview(post)
	if err != nil {
		zap.L().Error("create post with review failed", zap.Error(err))
		return nil, err
	}

	// 创建帖子审核任务
	reviewDao := mysql.NewPostReviewDao(ctx)
	task := &models.PostReviewTask{
		PostID:  postID,
		UserID:  user.UserId,
		Content: req.Content,
		Status:  "pending",
	}

	err = reviewDao.CreateReviewTask(task)
	if err != nil {
		zap.L().Error("create review task failed", zap.Error(err))
		return nil, err
	}

	// 保存到Redis
	err = redisCache.SaveReviewTask(task)
	if err != nil {
		zap.L().Error("save review task to redis failed", zap.Error(err))
		// 不返回错误，继续执行
	}

	// 创建审核日志
	log := &models.PostReviewLog{
		PostID:       postID,
		UserID:       user.UserId,
		Action:       "submit",
		OldStatus:    "",
		NewStatus:    "pending",
		ReviewResult: "帖子已提交审核",
	}

	err = reviewDao.CreateReviewLog(log)
	if err != nil {
		zap.L().Error("create review log failed", zap.Error(err))
		// 不返回错误，继续执行
	}

	// 异步发送到n8n进行审核
	go s.sendToN8nForReview(task)

	zap.L().Info("post submitted for review",
		zap.Int64("post_id", postID),
		zap.Int64("user_id", user.UserId))

	return &types.PostReviewResponse{
		PostID:       postID,
		ReviewStatus: "pending",
		Message:      "帖子已提交审核，请耐心等待",
	}, nil
}

// GetReviewStatus 获取帖子审核状态
func (s *ReviewSrv) GetReviewStatus(ctx context.Context, postID int64) (*types.ReviewStatusResponse, error) {
	// 首先尝试从Redis获取
	status, result, err := redisCache.GetReviewStatus(postID)
	if err == nil && status != "" {
		// Redis中有数据，直接返回
		return s.buildReviewStatusResponse(postID, status, result)
	}

	// Redis中没有数据，从数据库获取
	postDao := mysql.NewPostDao(ctx)
	post, err := postDao.GetPostWithReviewStatus(postID)
	if err != nil {
		zap.L().Error("get post review status failed", zap.Error(err))
		return nil, err
	}

	if post == nil {
		return &types.ReviewStatusResponse{
			PostID:  postID,
			Message: "帖子不存在",
		}, nil
	}

	// 解析审核结果
	var reviewResult *types.ReviewResult
	if post.ReviewResult != "" {
		err = json.Unmarshal([]byte(post.ReviewResult), &reviewResult)
		if err != nil {
			zap.L().Error("unmarshal review result failed", zap.Error(err))
		}
	}

	return s.buildReviewStatusResponse(postID, post.ReviewStatus, reviewResult)
}

// buildReviewStatusResponse 构建审核状态响应
func (s *ReviewSrv) buildReviewStatusResponse(postID int64, status string, result *types.ReviewResult) (*types.ReviewStatusResponse, error) {
	response := &types.ReviewStatusResponse{
		PostID:       postID,
		ReviewStatus: status,
	}

	if result != nil {
		response.ReviewScore = result.Score
		response.ReviewReason = result.Reason
		response.ReviewSuggestions = result.Suggestions
		response.ReviewTags = result.Tags
	}

	switch status {
	case "pending":
		response.Message = "帖子正在审核中，请耐心等待"
	case "approved":
		response.Message = "帖子审核通过"
	case "rejected":
		response.Message = "帖子审核未通过"
	default:
		response.Message = "未知审核状态"
	}

	return response, nil
}

// ProcessN8nCallback 处理n8n回调
func (s *ReviewSrv) ProcessN8nCallback(ctx context.Context, callback *types.ReviewCallbackRequest) (*types.ReviewCallbackResponse, error) {
	zap.L().Info("processing n8n callback",
		zap.Int64("post_id", callback.PostID),
		zap.String("status", callback.Status))

	// 构建审核结果
	reviewResult := &types.ReviewResult{
		Status:      callback.Status,
		Score:       callback.Score,
		Reason:      callback.Reason,
		Suggestions: callback.Suggestions,
		Tags:        callback.Tags,
		Details:     callback.Details,
	}

	// 更新数据库
	reviewDao := mysql.NewPostReviewDao(ctx)
	err := reviewDao.UpdatePostReviewStatus(callback.PostID, reviewResult)
	if err != nil {
		zap.L().Error("update post review status failed", zap.Error(err))
		return &types.ReviewCallbackResponse{
			Success: false,
			Message: "更新帖子审核状态失败",
		}, err
	}

	// 更新Redis缓存
	err = redisCache.SaveReviewStatus(callback.PostID, callback.Status, reviewResult)
	if err != nil {
		zap.L().Error("save review status to redis failed", zap.Error(err))
		// 不返回错误，继续执行
	}

	// 更新审核任务状态
	task, err := reviewDao.GetReviewTaskByPostID(callback.PostID)
	if err == nil && task != nil {
		status := "completed"
		if callback.Status == "rejected" {
			status = "failed"
		}
		err = reviewDao.UpdateReviewTaskStatus(task.ID, status, "")
		if err != nil {
			zap.L().Error("update review task status failed", zap.Error(err))
		}
	}

	// 创建审核日志
	log := &models.PostReviewLog{
		PostID:       callback.PostID,
		UserID:       callback.UserID,
		Action:       "review",
		OldStatus:    "pending",
		NewStatus:    callback.Status,
		ReviewResult: callback.Details,
	}

	err = reviewDao.CreateReviewLog(log)
	if err != nil {
		zap.L().Error("create review log failed", zap.Error(err))
		// 不返回错误，继续执行
	}

	// 发送通知给用户（这里可以集成通知系统）
	go s.sendReviewNotification(callback.PostID, callback.UserID, callback.Status, callback.Reason)

	zap.L().Info("n8n callback processed successfully",
		zap.Int64("post_id", callback.PostID),
		zap.String("status", callback.Status))

	return &types.ReviewCallbackResponse{
		Success: true,
		Message: "审核结果处理成功",
	}, nil
}

// sendToN8nForReview 发送审核请求到n8n
func (s *ReviewSrv) sendToN8nForReview(task *models.PostReviewTask) {
	// 构建发送给n8n的请求
	request := &types.ReviewRequest{
		PostID:    task.PostID,
		UserID:    task.UserID,
		Content:   task.Content,
		Timestamp: time.Now().Unix(),
	}

	// 序列化请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		zap.L().Error("marshal review request failed", zap.Error(err))
		return
	}

	// 发送HTTP请求到n8n
	// 从配置获取n8n URL
	reviewConfig := settings.DefaultReviewConfig()
	n8nURL := reviewConfig.N8nURL

	client := &http.Client{
		Timeout: time.Duration(reviewConfig.N8nTimeout) * time.Second,
	}

	resp, err := client.Post(n8nURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Error("send review request to n8n failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		zap.L().Error("n8n returned non-OK status",
			zap.Int("status", resp.StatusCode),
			zap.Int64("post_id", task.PostID))
		return
	}

	zap.L().Info("review request sent to n8n successfully",
		zap.Int64("post_id", task.PostID))
}

// sendReviewNotification 发送审核结果通知
func (s *ReviewSrv) sendReviewNotification(postID, userID int64, status, reason string) {
	// 这里可以集成通知系统，如邮件、短信、推送等
	zap.L().Info("sending review notification",
		zap.Int64("post_id", postID),
		zap.Int64("user_id", userID),
		zap.String("status", status))

	// 示例：发送邮件通知
	if status == "approved" {
		// 发送审核通过通知
		zap.L().Info("post approved notification sent", zap.Int64("post_id", postID))
	} else if status == "rejected" {
		// 发送审核拒绝通知
		zap.L().Info("post rejected notification sent",
			zap.Int64("post_id", postID),
			zap.String("reason", reason))
	}
}

// GetReviewStatistics 获取审核统计信息
func (s *ReviewSrv) GetReviewStatistics(ctx context.Context) (map[string]int64, error) {
	reviewDao := mysql.NewPostReviewDao(ctx)
	return reviewDao.GetReviewStatistics()
}

// GetPostsByReviewStatus 根据审核状态获取帖子列表
func (s *ReviewSrv) GetPostsByReviewStatus(ctx context.Context, status string, pageNum, pageSize int64) ([]*types.PostListResp, error) {
	postDao := mysql.NewPostDao(ctx)
	userDao := mysql.NewUserDao(ctx)
	communityDao := mysql.NewCommunityDao(ctx)

	posts, err := postDao.GetPostsByReviewStatus(status, pageNum, pageSize)
	if err != nil {
		zap.L().Error("get posts by review status failed", zap.Error(err))
		return nil, err
	}

	resp := make([]*types.PostListResp, 0, len(posts))
	for _, post := range posts {
		// 获取作者信息
		user, err := userDao.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("get user by id failed", zap.Error(err))
			continue
		}

		// 获取社区信息
		community, err := communityDao.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("get community by id failed", zap.Error(err))
			continue
		}

		postResp := &types.PostListResp{
			AuthorName: user.UserName,
			Post:       post,
			Community:  community,
		}

		resp = append(resp, postResp)
	}

	return resp, nil
}

// RetryFailedReview 重试失败的审核任务
func (s *ReviewSrv) RetryFailedReview(ctx context.Context, taskID int64) error {
	reviewDao := mysql.NewPostReviewDao(ctx)
	return reviewDao.RetryFailedReviewTask(taskID)
}
