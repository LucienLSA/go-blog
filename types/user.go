package types

// 定义请求的参数结构体

// 用户注册
type UserSignUpReq struct {
	Age        uint8  `form:"age" json:"age" binding:"gte=1,lte=130"`                             // 注册人年龄
	Gender     string `form:"gender" json:"gender" binding:"required,oneof=男 女 未知"`               // 注册人性别
	UserName   string `form:"user_name" json:"user_name" binding:"required"`                      // 注册人用户名
	Password   string `form:"password" json:"password" binding:"required"`                        // 注册人密码
	RePassword string `form:"re_password" json:"re_password" binding:"required,eqfield=Password"` // 注册人重复密码
	Email      string `form:"email" json:"email" binding:"omitempty"`                             // 注册人邮箱
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
	Avatar string `form:"avatar" json:"avatar" binding:"omitempty"` // 注册人默认头像
}
