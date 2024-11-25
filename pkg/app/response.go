package app

import (
	"blog-service/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// 返回统一的json格式数据
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"data": data,
		"msg":  e.GetMsg(errCode),
	})
}
