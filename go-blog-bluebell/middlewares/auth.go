package middlewares

import (
	"errors"
	"strings"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/pkg/jwt"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CtxUserKey = "user" // 用户信息的上下文key
	// 不需要验证token的路由
	SendEmailPath     = "/api/v2/email/send"
	ValidEmailPath    = "/api/v2/email/valid"
	SendEmailCodePath = "/api/v2/email/code/send"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 对于邮箱验证相关的路由，跳过token验证
		path := c.Request.URL.Path
		if path == SendEmailPath || path == ValidEmailPath || path == SendEmailCodePath {
			c.Next()
			return
		}

		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在请求头Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zap.L().Error("Request.Header.Get Authorization failed", zap.Error(errors.New("请求头中auth为空")))
			e.ResponseError(c, e.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			zap.L().Error("Request.Header.Get Authorization failed", zap.Error(errors.New("请求头中auth格式有误")))
			e.ResponseError(c, e.CodeNeedLogin)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("jwt.ParseToken failed", zap.Error(errors.New("无效的Token")))
			e.ResponseError(c, e.CodeTokenInvalid)
			c.Abort()
			return
		}
		// 从redis中获取token 并比较判断当前登录解析得到的token
		token, err := redisCache.GetJwtToken(mc.Username)
		// token不存在 需要重新登录
		if err == redisCache.ErrNotExistToken {
			e.ResponseError(c, e.CodeNeedLogin)
			c.Abort()
			return
		}
		// 	如果不一致，则说明在另一端登录
		if parts[1] != token {
			e.ResponseError(c, e.CodeLimitLogin)
			c.Abort()
			return
		}

		// 将用户信息保存到上下文中
		user := &types.User{
			UserID:   mc.UserID,
			UserName: mc.Username,
		}
		c.Set(CtxUserKey, user)

		// 设置用户信息到请求上下文
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{UserId: mc.UserID}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}
