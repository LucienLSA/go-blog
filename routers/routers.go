package routers

import (
	"github.com/LucienLSA/go-blog/controller"
	_ "github.com/LucienLSA/go-blog/docs"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin设置发布模式
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// 使用自定义的gin中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 导入swag接口文档
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
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
		// 参数动态获取帖子列表
		v1.GET("/searchposts", controller.SearchPostListHandler)
		// 根据社区查询帖子列表 整合到上一个hander中
		// v1.GET("/communitysearchposts", controller.CommunityPostListHandler)
		// 帖子投票
		v1.POST("/vote", controller.PostVoteHandler)
	}
	return r
}
