package routers

import (
	"blog-service/api"
	v1 "blog-service/api/v1"
	"blog-service/middleware"
	"blog-service/pkg/upload"
	"net/http"

	_ "blog-service/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
	// r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	article := v1.NewArticle()
	tag := v1.NewTag()
	author := api.NewAuthor()
	// 先只通过手动在数据库创建
	// r.POST("/author/register", author.RegisterAuthor) // 注册用户
	r.GET("/author/login", author.LoginAuthor) // 登录用户
	apiv1 := r.Group("/api/v1")
	// fmt.Println("JWT sercret:", global.AppSetting.JwtSecret)
	apiv1.Use(middleware.JWT()) // JWT中间件
	{
		apiv1.POST("/tags", tag.CreateTags)       // 新增标签
		apiv1.DELETE("/tags/:id", tag.DeleteTags) // 删除标签
		apiv1.PUT("/tags/:id", tag.UpdateTags)    // 更新标签
		// apiv1.PATCH("/tags/:id/state", tag.UpdateTag) // 更新标签状态
		apiv1.GET("/tags", tag.GetTags) // 获取标签列表

		apiv1.POST("/articles", article.CreateArticles)       // 新增文章
		apiv1.DELETE("/articles/:id", article.DeleteArticles) // 删除文章
		apiv1.PUT("/articles/:id", article.UpdateArticles)    // 更新文章
		// apiv1.PATCH("/articles/:id/state", article.UpdateArticle) // 更新文章状态
		apiv1.GET("/articles/:id", article.GetArticle) // 获取指定文章详情
		apiv1.GET("/articles", article.ListArticles)   // 获取文章列表
	}
	return r
}
