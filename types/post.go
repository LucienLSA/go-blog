package types

// 根据请求的参数查询方式
const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 投票数据
type PostVoteDataReq struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id,string" binding:"required"`  // 投票帖子id
	Kind   int8  `json:"kind,string" binding:"oneof=1 0 -1"` // 投票帖子赞成(1)\反对(-1)\取消投票(0)
}

// 帖子分类查询 id
type PostIdReq struct {
	PostID int64 `json:"post_id,string" form:"community_id" binding:"required"` // 帖子id
}

// 帖子发布参数
type PostCreateReq struct {
	AuthorID    int64  `json:"author_id" db:"author_id"`                          // 发布者id
	CommunityID int64  `json:"community_id" db:"community_id" binding:"required"` // 社区id
	Title       string `json:"title" db:"title" binding:"required"`               // 帖子标题
	Content     string `json:"content" db:"content" binding:"required"`           // 帖子内容
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

// // 获取帖子列表query string参数 和 社区id
// type CommunityPostListReq struct {
// 	*ParamPostList
// 	CommunityID int64 `json:"community_id" form:"community_id"`
// }
