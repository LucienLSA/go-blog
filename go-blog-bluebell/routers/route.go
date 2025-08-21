package routers

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/controller"
	_ "github.com/LucienLSA/go-blog/docs"
	logging "github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/middlewares"
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

	// 处理异常
	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	// 使用自定义的ginlogger中间件
	r.Use(logging.GinLogger(), logging.GinRecovery(true))
	// // 存入cookie
	// store := cookie.NewStore([]byte("something-very-secret"))
	// r.Use(sessions.Sessions("mysession", store))

	// 使用跨域中间件
	r.Use(middlewares.Cors())

	// 加载静态文件和html
	r.LoadHTMLFiles("./templates/index.html", "./templates/oauth_success.html", "./templates/oauth_error.html")
	r.Static("/static", "./static")

	// 代入根目录的html
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	// Static test pages for OAuth
	r.GET("/oauth/github/callback", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "oauth_success.html", nil)
	})
	r.GET("/oauth/github/error", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "oauth_error.html", nil)
	})

	// 导入swag接口文档
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	v2 := r.Group("/api/v2")

	userRoute := v2.Group("/user")
	{
		// // 注册
		userRoute.POST("/signup", controller.SignUpHandler())

		// // 登录
		userRoute.POST("/login", controller.LoginHandler())
		userRoute.POST("/login/emailcode", controller.SendEmailCodeHandler())
		userRoute.POST("/login/email", controller.LoginEmailHandler())

		userAuthRoute := userRoute.Use(middlewares.JWTAuthMiddleware())
		// // 用户头像上传
		userAuthRoute.POST("/avatar", controller.UploadAvatarHandler())

		// // 用户信息更新
		userAuthRoute.PUT("/update", controller.UpdateHandler())

		// // 用户发送邮箱
		userAuthRoute.POST("/sendEmail", controller.SendEmailHandler())

		// // 用户验证邮箱
		userAuthRoute.POST("/validEmail", controller.ValidEmailHandler())
		userAuthRoute.GET("/validEmail", controller.ValidEmailHandler())

		// // 用户登出
		userAuthRoute.DELETE("/logout", controller.LogoutHandler())
	}
	communityRoute := v2.Group("/community")
	{
		// // 获取社区信息
		communityRoute.GET("/show", controller.CommunityListHandler())
		communityRoute.GET("/show/:community_id", controller.CommunityDetailHandler())
		// 创建社区
		communityRoute.Use(middlewares.JWTAuthMiddleware()).POST("/create", controller.CreateCommunityHandler())
	}
	postsRoute := v2.Group("/posts")
	{
		// // 查看帖子
		postsRoute.GET("/showAll", controller.GetPostListHandler())
		postsAuthRoute := postsRoute.Use(middlewares.JWTAuthMiddleware())
		// // 帖子查询新版
		// // 参数动态获取帖子列表
		postsAuthRoute.GET("/search", controller.SearchPostListHandler())

		// // 根据社区查询帖子列表 整合到上一个hander中
		// // v1.GET("/communitysearchposts", controller.CommunityPostListHandler)

		// // 令牌桶填充速率2s， 容量1
		postsAuthRoute.GET("/:post_id(\\d{1,18})", middlewares.RateLimitMiddleware(2*time.Second, 1), controller.GetPostDetailHandler())

		// // 发布帖子
		postsAuthRoute.POST("/post", controller.CreatePostHandler())

		// // 帖子投票
		postsAuthRoute.POST("/vote", controller.PostVoteHandler())

		// 帖子审核
		// 提交帖子审核
		postsRoute.POST("/review", controller.SubmitPostForReviewHandler())
		// 查询帖子审核状态
		postsRoute.GET("/review/:post_id", controller.GetReviewStatusHandler())
		// n8n回调
		postsRoute.POST("/review/callback", controller.N8nReviewCallbackHandler())
		// 获取审核统计
		postsRoute.GET("/review/statistics", controller.GetReviewStatisticsHandler())
		// 获取审核列表
		postsRoute.GET("/review/list", controller.GetPostsByReviewStatusHandler())
	}

	// GitHub热点数据路由
	githubRoute := v2.Group("/github")
	{
		// 获取GitHub热点数据
		githubRoute.GET("/trending", controller.GetGitHubTrendingHandler())
		// 获取所有GitHub热点数据
		githubRoute.GET("/trending/all", controller.GetAllGitHubTrendingHandler())
		// 手动刷新GitHub热点数据
		githubRoute.POST("/trending/refresh", controller.RefreshGitHubTrendingHandler())
	}

	// OAuth GitHub 登录路由
	oauthGitHub := v2.Group("/oauth/github")
	{
		oauthGitHub.GET("/login", controller.GitHubLoginHandler())
		oauthGitHub.GET("/callback", controller.GitHubCallbackHandler())
	}

	// 注册pprof路由 服务型性能分析
	// pprof.Register(r)

	return r
}

// 未找到资源
func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "404",
	})
}
