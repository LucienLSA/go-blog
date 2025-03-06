package route

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/controller"
	_ "github.com/LucienLSA/go-blog/docs"
	"github.com/LucienLSA/go-blog/logger"
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
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// // 存入cookie
	// store := cookie.NewStore([]byte("something-very-secret"))
	// r.Use(sessions.Sessions("mysession", store))

	// 使用跨域中间件
	r.Use(middlewares.Cors())

	// 加载静态文件和html
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	// 代入根目录的html
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// 导入swag接口文档
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/user/signup", controller.SignUpHandler)

	// 登录
	v1.POST("/user/login", controller.LoginHandler)
	v1.POST("/user/login/email", controller.LoginEmailHandler)

	// 获取社区信息
	v1.GET("/community/show", controller.CommunityHandler)
	v1.GET("/community/show/:community_id", controller.CommunityDetailHandler)

	// 查看帖子
	v1.GET("/posts/showAll", controller.GetPostListHandler)

	// JWT中间件认证
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// 帖子查询新版
		// 参数动态获取帖子列表
		v1.GET("/posts/search", controller.SearchPostListHandler)

		// 根据社区查询帖子列表 整合到上一个hander中
		// v1.GET("/communitysearchposts", controller.CommunityPostListHandler)

		// 令牌桶填充速率2s， 容量1
		v1.GET("/posts/:post_id", middlewares.RateLimitMiddleware(2*time.Second, 1), controller.GetPostDetailHandler)

		// 发布帖子
		v1.POST("/posts/post", controller.CreatePostHandler)

		// 帖子投票
		v1.POST("/posts/vote", controller.PostVoteHandler)

		// 用户头像上传
		v1.POST("/user/avatar", controller.UploadAvatarHandler)

		// 用户信息更新
		v1.POST("/user/update", controller.UpdateHandler)

		// 用户发送邮箱
		v1.POST("/user/sendEmail", controller.SendEmailHandler)

		// 用户验证邮箱
		v1.GET("/user/validEmail", controller.ValidEmailHandler)
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
