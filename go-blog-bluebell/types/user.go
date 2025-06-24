package types

// User 用户基本信息
type User struct {
	UserID   int64  `json:"user_id" db:"user_id"`     // 用户ID
	UserName string `json:"user_name" db:"user_name"` // 用户名
	Gender   string `json:"gender" db:"gender"`       // 用户性别
	Age      uint8  `json:"age" db:"age"`             // 用户年龄
	Email    string `json:"email" db:"email"`         // 用户邮箱
	Avatar   string `json:"avatar" db:"avatar"`       // 用户头像
}

// 定义请求的参数结构体

// 用户注册
type UserSignUpReq struct {
	Age        uint8  `form:"age" json:"age" binding:"gte=1,lte=130"`                             // 注册人年龄
	Gender     string `form:"gender" json:"gender" binding:"required,oneof=男 女 未知"`               // 注册人性别
	UserName   string `form:"user_name" json:"user_name" binding:"required"`                      // 注册人用户名
	Password   string `form:"password" json:"password" binding:"required"`                        // 注册人密码
	RePassword string `form:"re_password" json:"re_password" binding:"required,eqfield=Password"` // 注册人重复密码
	Email      string `form:"email" json:"email" binding:"omitempty,email"`                       // 注册人邮箱
	// // 需要使用自定义校验方法checkDate做参数校验的字段Date
	// Date string `json:"date" binding:"required,datetime=2006-01-02,checkDate"`
	Avatar string `form:"avatar" json:"avatar" binding:"omitempty"` // 注册人默认头像
}

// 用户信息更新
type UserUpdateReq struct {
	Age         uint8  `form:"age" json:"age" binding:"omitempty,gte=1,lte=130"`      // 用户年龄
	Gender      string `form:"gender" json:"gender" binding:"omitempty,oneof=男 女 未知"` // 用户性别
	UserName    string `form:"user_name" json:"user_name" binding:"omitempty"`        // 用户名
	Password    string `form:"password" json:"password" binding:"omitempty"`          // 当前密码
	NewPassword string `form:"new_password" json:"new_password" binding:"omitempty"`  // 新密码
	Email       string `form:"email" json:"email" binding:"omitempty,email"`          // 用户邮箱
	Avatar      string `form:"avatar" json:"avatar" binding:"omitempty"`              // 用户头像
}

// 用户登录
type UserLoginReq struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"` // 登录人姓名
	Password string `json:"password" form:"password" binding:"required"`   // 登录人密码
}

// 用户登录返回
type UserLoginResp struct {
	UserID   int64  `json:"user_id" form:"user_id"`     // 用户id
	UserName string `json:"user_name" form:"user_name"` // 登录人姓名
	Token    string `json:"token" form:"token"`         // 用户token
}

// 用户邮箱发送验证码绑定和解绑
type UserSendEmailReq struct {
	OperationType int    `form:"operation_type" json:"operation_type" binding:"required,oneof=1 2"` // 邮箱操作类型 1-绑定 2-解绑
	Email         string `form:"email" json:"email" binding:"required,email"`                       // 用户邮箱
}

// 用户验证邮箱(绑定或解绑)
type UserVaildEmail struct {
	UserName      string `json:"user_name" binding:"-"` // 禁止自动绑定
	UserID        int64  `json:"user_id" binding:"-"`   // 禁止自动绑定
	Email         string `form:"email" json:"email" binding:"required,email"`
	Code          string `form:"code" json:"code" binding:"required"`
	OperationType int    `form:"operation_type" json:"operation_type" binding:"required,oneof=1 2"`
}

// 用户发送邮箱登录验证码
type UserSendEmailCodeReq struct {
	OperationType int    `form:"operation_type" json:"operation_type" binding:"required,oneof=3"` // 邮箱操作类型 3-登录
	UserEmail     string `form:"email" json:"email" binding:"email,required"`                     // 用户邮箱
}

// 用户邮箱登录
type UserLoginEmailReq struct {
	UserEmail string `form:"email" json:"email" binding:"email,required"` // 用户邮箱
	Code      string `form:"code" json:"code" binding:"required"`         // 邮箱收到的验证码
}

// 用户头像上传
type UserAvatar struct {
	UserID   int64  `form:"user_id" json:"user_id"`     // 用户ID
	UserName string `form:"user_name" json:"user_name"` // 用户名称
	Avatar   string `form:"avatar" json:"avatar"`       // 用户密头像
}
