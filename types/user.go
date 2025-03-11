package types

// 定义请求的参数结构体

// 用户注册
type UserSignUpReq struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`                     // 注册人年龄
	Gender     string `json:"gender" binding:"required,oneof=男 女 未知"`          // 注册人性别
	Username   string `json:"username" binding:"required"`                     // 注册人用户名
	Password   string `json:"password" binding:"required"`                     // 注册人密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 注册人重复密码
	Email      string `json:"email" binding:"omitempty"`                       // 注册人邮箱
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
	Avatar string `json:"avatar" binding:"omitempty"` // 注册人默认头像
}
