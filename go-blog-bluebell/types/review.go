package types

// PostReviewRequest 帖子审核请求
type PostReviewRequest struct {
	PostID      int64  `json:"post_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommunityID int64  `json:"community_id" binding:"required"`
}

// PostReviewResponse 帖子审核响应
type PostReviewResponse struct {
	PostID       int64  `json:"post_id"`
	ReviewStatus string `json:"review_status"`
	Message      string `json:"message"`
}

// ReviewStatusRequest 获取审核状态请求
type ReviewStatusRequest struct {
	PostID int64 `json:"post_id" binding:"required"`
}

// ReviewStatusResponse 审核状态响应
type ReviewStatusResponse struct {
	PostID            int64    `json:"post_id"`
	ReviewStatus      string   `json:"review_status"`
	ReviewScore       float64  `json:"review_score"`
	ReviewReason      string   `json:"review_reason"`
	ReviewSuggestions string   `json:"review_suggestions"`
	ReviewTags        []string `json:"review_tags"`
	ReviewedAt        string   `json:"reviewed_at"`
	Message           string   `json:"message"`
}

// ReviewResult n8n审核结果
type ReviewResult struct {
	Status      string   `json:"status"`      // approved, rejected
	Score       float64  `json:"score"`       // 审核评分 0-100
	Reason      string   `json:"reason"`      // 拒绝原因
	Suggestions string   `json:"suggestions"` // 修改建议
	Tags        []string `json:"tags"`        // 内容标签
	Details     string   `json:"details"`     // 详细结果
}

// ReviewRequest 发送给n8n的审核请求
type ReviewRequest struct {
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ReviewCallback n8n回调数据
type ReviewCallback struct {
	PostID       int64        `json:"post_id"`
	UserID       int64        `json:"user_id"`
	ReviewResult ReviewResult `json:"review_result"`
	ProcessedAt  int64        `json:"processed_at"`
}

// ReviewCallbackRequest n8n回调请求
type ReviewCallbackRequest struct {
	PostID      int64    `json:"post_id" binding:"required"`
	UserID      int64    `json:"user_id" binding:"required"`
	Status      string   `json:"status" binding:"required"`
	Score       float64  `json:"score"`
	Reason      string   `json:"reason"`
	Suggestions string   `json:"suggestions"`
	Tags        []string `json:"tags"`
	Details     string   `json:"details"`
	ProcessedAt int64    `json:"processed_at"`
}

// ReviewCallbackResponse n8n回调响应
type ReviewCallbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
