package response

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist

	CodeInvalidPassword
	CodeServerBusy

	CodeEmailNotExist
	CodeEmailExist

	CodeNeedLogin = 2000 + iota
	CodeTokenInvalid
	CodeLimitLogin

	CodeUploadFile = 3000 + iota
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务出错",
	CodeNeedLogin:       "需要登录",
	CodeTokenInvalid:    "无效的Token",
	CodeLimitLogin:      "已在另一台设备登录",
	CodeUploadFile:      "上传文件失败",
	CodeEmailNotExist:   "邮箱不存在",
	CodeEmailExist:      "邮箱已存在",
}
