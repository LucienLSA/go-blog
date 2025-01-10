package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10001, // 错误码
	"msg": xx, // 提示信息
	"data": {}, // 数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 返回错误
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.getMsg(),
		Data: nil,
	})
}

// 返回错误带信息
func ResponseErrorMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// 返回成功带数据
func ResponseSuccessData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.getMsg(),
		Data: data,
	})
}
