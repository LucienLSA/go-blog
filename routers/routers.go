package routers

import (
	"net/http"

	"github.com/LucienLSA/go-blog/controller"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin设置发布模式
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	// 注册
	r.POST("/signup", controller.SignUpHandler)
	// 登录
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middleware.JWTAuthMiddleware(), func(ctx *gin.Context) {
		// 登录用户 ，判断请求头中存在有效的JWT
		ctx.String(http.StatusOK, "pong")
	})

	return r
}
