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
	Gender     string `json:"gender" binding:"required,oneof=男 女 未知"`          // 注册人性别
	Username   string `json:"username" binding:"required"`                     // 注册人用户名
	Password   string `json:"password" binding:"required"`                     // 注册人密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 注册人重复密码
	Email      string `json:"email"`                                           // 注册人邮箱
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
	Avatar string `json:"avatar"` // 注册人默认头像
}
type ParamUpdate struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`                     // 用户年龄
	Gender     string `json:"gender" binding:"required,oneof=男 女 未知"`          // 用户性别
	Username   string `json:"username" binding:"required"`                     // 用户名
	Password   string `json:"password" binding:"required"`                     // 用户密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 用户重复密码
	Email      string `json:"email"`                                           // 用户邮箱
	Avatar     string `json:"avatar"`                                          // 用户头像
}

// 用户登录
type ParamLogin struct {
	Username string `json:"username" form:"username" binding:"required"`  // 登录人姓名
	Password string ` json:"password" form:"password" binding:"required"` // 登录人密码
}

// 用户头像上传
type ParamAvatar struct {
	UserID   int64  `form:"user_id" json:"user_id"`   // 用户ID
	UserName string `form:"username" json:"username"` // 用户名称
	Avatar   string `form:"avatar" json:"avatar"`     // 用户密头像
}

// 用户邮箱发送验证码
type ParamSendEmail struct {
	OperationType int    `form:"operation_type" json:"operation_type" binding:"required"`
	Email         string `form:"email" json:"email" binding:"required, email"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id,string" binding:"required"`  // 投票帖子id
	Kind   int8  `json:"kind,string" binding:"oneof=1 0 -1"` // 投票帖子赞成(1)\反对(-1)\取消投票(0)
}

// 社区分类查询 id
type ParamCommunityId struct {
	CommunityID int64 `json:"community_id" form:"community_id" example:"2" binding:"required"` // 社区id
}

// 帖子分类查询 id
type ParamPostId struct {
	PostID int64 `json:"post_id,string" form:"community_id" binding:"required"` // 帖子id
}

// 帖子发布参数
type ParamPostCreate struct {
	CommunityID int64  `json:"community_id" db:"community_id" binding:"required"` // 社区id
	Title       string `json:"title" db:"title" binding:"required"`               // 帖子标题
	Content     string `json:"content" db:"content" binding:"required"`           // 帖子内容
}

// 帖子分页查询
type ParamPostSearch struct {
	PageNum  int64 `json:"page_num" form:"page_num" example:"1"`   // 获取帖子列表的页码
	PageSize int64 `json:"page_size" form:"page_size" example:"5"` // 获取帖子列表的数量
}

// 帖子列表query string参数 社区id 整合后的
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
