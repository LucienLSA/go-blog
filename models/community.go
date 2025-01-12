package models

import "time"

type Community struct {
	CommunityID   int64  `json:"community_id" db:"community_id"`     // 社区id
	CommunityName string `json:"community_name" db:"community_name"` // 社区名称
}

type CommunityDetail struct {
	CommunityID   int64     `json:"community_id" db:"community_id"`           // 社区id
	CommunityName string    `json:"community_name" db:"community_name"`       // 社区名称
	Introduction  string    `json:"introduction,omitempty" db:"introduction"` // 社区介绍
	CreateTime    time.Time `json:"create_time" db:"create_time"`             // 社区创建时间
	UpdateTime    time.Time `json:"update_time" db:"update_time"`             // 社区更新时间
}
