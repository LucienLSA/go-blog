package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PostID      uint   `gorm:"post_id"`      // 帖子id
	CommunityID uint   `gorm:"community_id"` // 社区id
	AuthorID    uint   ` gorm:"author_id"`   // 帖子作者id
	Status      int32  ` gorm:"status"`      // 帖子状态
	Title       string ` gorm:"title"`       // 帖子标题
	Content     string ` gorm:"content"`     // 帖子内容
}

// type ApiPostDetail struct {
// 	AuthorName       string             `json:"author_name"`    // 帖子作者用户名
// 	VoteAgreeNum     int64              `json:"vote_agree_num"` // 投票赞成数
// 	*Post            `json:"post"`      // 嵌入帖子信息
// 	*CommunityDetail `json:"community"` // 嵌入社区信息
// }
