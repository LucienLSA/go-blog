package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Status      int32  ` gorm:"status"`      // 帖子状态 0:帖子删除；1：帖子活跃
	PostID      int64  `gorm:"post_id"`      // 帖子id
	CommunityID int64  `gorm:"community_id"` // 社区id
	AuthorID    int64  ` gorm:"author_id"`   // 帖子作者id
	Title       string ` gorm:"title"`       // 帖子标题
	Content     string ` gorm:"content"`     // 帖子内容
	// 审核相关字段
	ReviewStatus      string  `gorm:"review_status"`       // 审核状态: pending, approved, rejected
	ReviewResult      string  `gorm:"review_result"`       // 审核结果详情
	ReviewedAt        *string `gorm:"reviewed_at"`         // 审核时间
	ReviewScore       float64 `gorm:"review_score"`        // 审核评分
	ReviewReason      string  `gorm:"review_reason"`       // 拒绝原因
	ReviewSuggestions string  `gorm:"review_suggestions"`  // 修改建议
	ReviewTags        string  `gorm:"review_tags"`         // 内容标签(JSON格式)
}

// type ApiPostDetail struct {
// 	AuthorName       string             `json:"author_name"`    // 帖子作者用户名
// 	VoteAgreeNum     int64              `json:"vote_agree_num"` // 投票赞成数
// 	*Post            `json:"post"`      // 嵌入帖子信息
// 	*CommunityDetail `json:"community"` // 嵌入社区信息
// }
