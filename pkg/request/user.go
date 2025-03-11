package request

import (
	"errors"

	"github.com/LucienLSA/go-blog/middlewares"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取到经过中间件的jwt auth鉴权后登录的用户信息 ，进行下一步处理
func GetLoginUserID(c *gin.Context) (userID uint, err error) {
	uid, ok := c.Get(middlewares.CtxtUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	if userID != uid {
		err = ErrorUserNotLogin
		return
	}
	return
}
