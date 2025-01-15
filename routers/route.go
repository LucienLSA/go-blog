package route

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/controller"
	_ "github.com/LucienLSA/go-blog/docs"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/middlewares"
	"github.com/gin-contrib/pprof"
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
	//处理异常
	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)
	// 使用自定义的ginlogger中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
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

	// 注册业务路由
	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)
	// 获取社区信息
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:community_id", controller.CommunityDetailHandler)
	// 查看帖子
	v1.GET("/posts", controller.GetPostListHandler)
	// 帖子查询新版
	// 参数动态获取帖子列表
	v1.GET("/searchposts", controller.SearchPostListHandler)
	// 根据社区查询帖子列表 整合到上一个hander中
	// v1.GET("/communitysearchposts", controller.CommunityPostListHandler)
	// 令牌桶填充速率2s， 容量1
	v1.GET("/post/:post_id", middlewares.RateLimitMiddleware(2*time.Second, 1), controller.GetPostDetailHandler)
	// JWT中间件认证
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		// 帖子
		v1.POST("/post", controller.CreatePostHandler)

		// 帖子投票
		v1.POST("/vote", controller.PostVoteHandler)
	}
	// 注册pprof路由 服务型性能分析
	pprof.Register(r)

	return r
}

// 未找到资源
func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "404",
	})
}
