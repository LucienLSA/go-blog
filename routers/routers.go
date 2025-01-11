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
		// 社区信息
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:community_id", controller.CommunityDetailHandler)
		// 帖子
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.GET("/post/:post_id", controller.GetPostDetailHandler)
		// 帖子查询新版
		v1.GET("/searchposts", controller.SearchPostListHandler)
		// 帖子投票
		v1.POST("/vote", controller.PostVoteHandler)
	}
	return r
}
