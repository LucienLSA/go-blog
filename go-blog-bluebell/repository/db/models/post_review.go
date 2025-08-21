package models

import (
	"time"
)

// PostReviewTask 帖子审核任务
type PostReviewTask struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID       int64      `gorm:"post_id;not null" json:"post_id"`
	UserID       int64      `gorm:"user_id;not null" json:"user_id"`
	Content      string     `gorm:"content;type:text;not null" json:"content"`
	Status       string     `gorm:"status;default:'pending'" json:"status"` // pending, processing, completed, failed
	CreatedAt    time.Time  `gorm:"created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CompletedAt  *time.Time `gorm:"completed_at" json:"completed_at"`
	ErrorMessage string     `gorm:"error_message;type:text" json:"error_message"`
}

// PostReviewLog 帖子审核日志
type PostReviewLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID       int64     `gorm:"post_id;not null" json:"post_id"`
	UserID       int64     `gorm:"user_id;not null" json:"user_id"`
	Action       string    `gorm:"action;not null" json:"action"` // submit, approve, reject, update
	OldStatus    string    `gorm:"old_status" json:"old_status"`
	NewStatus    string    `gorm:"new_status" json:"new_status"`
	ReviewResult string    `gorm:"review_result;type:text" json:"review_result"`
	OperatorID   *int64    `gorm:"operator_id" json:"operator_id"`
	CreatedAt    time.Time `gorm:"created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}
