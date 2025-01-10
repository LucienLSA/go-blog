package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
}

// type PostList struct {
// 	PostID int64  `json:"post_id" db:"post_id"`
// 	Title  string `json:"title" db:"title" binding:"required"`
// 	Status int32  `json:"status" db:"status"`
// }

type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	*Post            `json:"post"`      // 嵌入帖子信息
	*CommunityDetail `json:"community"` // 嵌入社区信息
}
