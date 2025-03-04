package models

import "time"

// Notice 公告模型 存放公告和邮件模板
type Notice struct {
	Text       string    `json:"text" db:"text"`
	CreateTime time.Time `json:"create_time" db:"create_time"` // 帖子创建时间
	UpdateTime time.Time `json:"update_time" db:"update_time"` // 帖子更新时间
}
