package types

import "github.com/LucienLSA/go-blog/repository/db/models"

// 根据请求的参数查询方式
const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 帖子分类查询 id
type PostIdReq struct {
	PostID int64 `json:"post_id,string" form:"post_id" binding:"required"` // 帖子id
}

// 帖子发布参数
type PostCreateReq struct {
	PostID      int64  `json:"post_id" form:"post_id"`                                // 帖子id
	AuthorID    int64  `json:"author_id" form:"author_id"`                            // 发布者id
	CommunityID int64  `json:"community_id" dformb:"community_id" binding:"required"` // 社区id
	Title       string `json:"title" form:"title" binding:"required"`                 // 帖子标题
	Content     string `json:"content" form:"content" binding:"required"`             // 帖子内容
}

// 帖子分页查询
type PostSearchReq struct {
	PageNum  int64 `json:"page_num" form:"page_num" example:"1"`   // 获取帖子列表的页码
	PageSize int64 `json:"page_size" form:"page_size" example:"5"` // 获取帖子列表的数量
}

// 帖子列表query string参数 社区id 整合后的
type PostListReq struct {
	PageNum  int64  `json:"page_num" form:"page_num" example:"1"`   // 获取帖子列表的页码
	PageSize int64  `json:"page_size" form:"page_size" example:"5"` // 获取帖子列表的数量
	Order    string `json:"order" form:"order" example:"score"`     // 排序方式 按照时间或者投票分数

	CommunityID int64 `json:"community_id" form:"community_id"` // 社区id可以为空
}

type PostListResp struct {
	AuthorName              string                              `json:"author_name" form:"author_name"`       // 帖子作者用户名
	VoteAgreeNum            int64                               `json:"vote_agree_num" form:"vote_agree_num"` // 投票赞成数
	*models.Post            `json:"post" form:"post"`           // 嵌入帖子信息
	*models.CommunityDetail `json:"community" form:"community"` // 嵌入社区信息
}

// type ApiPostDetail struct {
// 	AuthorName       string             `json:"author_name"`    // 帖子作者用户名
// 	VoteAgreeNum     int64              `json:"vote_agree_num"` // 投票赞成数
// 	*Post            `json:"post"`      // 嵌入帖子信息
// 	*CommunityDetail `json:"community"` // 嵌入社区信息
// }

// // 获取帖子列表query string参数 和 社区id
// type CommunityPostListReq struct {
// 	*ParamPostList
// 	CommunityID int64 `json:"community_id" form:"community_id"`
// }
