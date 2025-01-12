package models

// 定义请求的参数结构体

// 根据请求的参数查询方式
const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 注册
type ParamSignUp struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`                     // 注册人年龄
	Gender     string `json:"gender" binding:"required,oneof=男 女"`             // 注册人性别
	Username   string `json:"username" binding:"required"`                     // 注册人用户名
	Password   string `json:"password" binding:"required"`                     // 注册人密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 注册人重复密码
	Email      string `json:"email" binding:"required,email"`                  // 注册人邮箱
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
}

// 登录
type ParamLogin struct {
	Username string `json:"username" binding:"required"` // 登录人姓名
	Password string `json:"password" binding:"required"` // 登录人密码
}

// 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id,string" binding:"required"`  // 投票帖子id
	Kind   int8  `json:"kind,string" binding:"oneof=1 0 -1"` // 投票帖子赞成(1)\反对(-1)\取消投票(0)
}

// 获取帖子列表query string参数 社区id 整合后的
type ParamPostList struct {
	PageNum  int64  `json:"page_num" form:"page_num" example:"1"`   // 获取帖子列表的页码
	PageSize int64  `json:"page_size" form:"page_size" example:"5"` // 获取帖子列表的数量
	Order    string `json:"order" form:"order" example:"score"`     // 排序方式 按照时间或者投票分数

	CommunityID int64 `json:"community_id" form:"community_id"` // 社区id可以为空
}

// // 获取帖子列表query string参数 和 社区id
// type ParamCommunityPostList struct {
// 	*ParamPostList
// 	CommunityID int64 `json:"community_id" form:"community_id"`
// }
