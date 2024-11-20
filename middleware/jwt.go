package middleware

import (
	"blog-service/pkg/e"
	"blog-service/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		// var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.ERROR_AUTH_TOKEN
			c.JSON(e.INVALID_PARAMS, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": "缺失token",
			})
			// 阻止后续中间件或程序执行
			c.Abort()
			return
		} else {
			claim, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				c.JSON(e.INVALID_PARAMS, gin.H{
					"code": code,
					"msg":  e.GetMsg(code),
					"data": "解析token失败",
				})
				c.Abort()
				return
			} else if time.Now().Unix() > claim.ExpiresAt { // Token has expired
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				c.JSON(e.INVALID_PARAMS, gin.H{
					"code": code,
					"msg":  e.GetMsg(code),
					"data": "token已过期",
				})

				c.Abort()
				return
			}
		}
		c.Next()

	}

}
