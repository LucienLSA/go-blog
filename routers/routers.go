package routers

import (
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

	v1 := r.Group("/api/v1")

	// 注册业务路由
	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	// JWT中间件认证
	v1.Use(middleware.JWTAuthMiddleware())

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
	}
	return r
}
