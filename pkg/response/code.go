package response

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeJWTheaderAuthEmpty = 2000 + iota
	CodeJWTheaderAuthError
	CodeJWTtokenInvalid
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "请求参数错误",
	CodeUserExist:          "用户名已存在",
	CodeUserNotExist:       "用户名不存在",
	CodeInvalidPassword:    "用户名或密码错误",
	CodeServerBusy:         "服务出错",
	CodeJWTheaderAuthEmpty: "请求头中auth为空",
	CodeJWTheaderAuthError: "请求头中auth格式有误",
	CodeJWTtokenInvalid:    "无效的Token",
}
