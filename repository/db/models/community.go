package models

import "gorm.io/gorm"

type Community struct {
	gorm.Model
	CommunityID   int64  ` gorm:"community_id"`  // 社区id
	CommunityName string `gorm:"community_name"` // 社区名称
}

// type CommunityDetail struct {
// 	CommunityID   int64     ` gorm:"community_id"`   // 社区id
// 	CommunityName string    ` gorm:"community_name"` // 社区名称
// 	Introduction  string    ` gorm:"introduction"`   // 社区介绍
// 	CreateTime    time.Time ` gorm:"create_time"`    // 社区创建时间
// 	UpdateTime    time.Time ` gorm:"update_time"`    // 社区更新时间
// }
