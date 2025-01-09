package request

import (
	"errors"

	"github.com/LucienLSA/go-blog/middleware"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取到经过中间件的jwt auth鉴权后登录的用户信息 ，进行下一步处理
func GetLoginUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(middleware.CtxtUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
