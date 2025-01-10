package models

// 定义请求的参数结构体
// 注册
type ParamSignUp struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Gender     string `json:"gender" binding:"required,oneof=男 女"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
}

// 登录
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id,string" binding:"required"`           // 帖子
	Kind   int8  `json:"kind,string" binding:"required,oneof=1 0 -1"` // 赞成(1)\反对(-1)\取消投票(0)
}
