package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`                       // 帖子id
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"` // 社区id
	AuthorID    int64     `json:"author_id,string" db:"author_id"`                   // 帖子作者id
	Status      int32     `json:"status" db:"status"`                                // 帖子状态
	CreateTime  time.Time `json:"create_time" db:"create_time"`                      // 帖子创建时间
	UpdateTime  time.Time `json:"update_time" db:"update_time"`                      // 帖子更新时间
	Title       string    `json:"title" db:"title" binding:"required"`               // 帖子标题
	Content     string    `json:"content" db:"content" binding:"required"`           // 帖子内容
}

// type PostList struct {
// 	PostID int64  `json:"post_id" db:"post_id"`
// 	Title  string `json:"title" db:"title" binding:"required"`
// 	Status int32  `json:"status" db:"status"`
// }

type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`    // 帖子作者用户名
	VoteAgreeNum     int64              `json:"vote_agree_num"` // 投票赞成数
	*Post            `json:"post"`      // 嵌入帖子信息
	*CommunityDetail `json:"community"` // 嵌入社区信息
}
