package routers

import (
	v1 "blog-service/api/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	// r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())
	r := gin.Default()
	article := v1.NewArticle()
	tag := v1.NewTag()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/tags", tag.CreateTag)            // 新增标签
		apiv1.DELETE("/tags/:id", tag.DeleteTag)      // 删除标签
		apiv1.PUT("/tags/:id", tag.UpdateTag)         // 修改标签
		apiv1.PATCH("/tags/:id/state", tag.UpdateTag) // 更新标签状态
		apiv1.GET("/tags", tag.ListTag)               // 获取标签列表

		apiv1.POST("/articles", article.CreateArticle)            // 新增文章
		apiv1.DELETE("/articles/:id", article.DeleteArticle)      // 删除文章
		apiv1.PUT("/articles/:id", article.UpdateArticle)         // 修改文章
		apiv1.PATCH("/articles/:id/state", article.UpdateArticle) // 更新文章状态
		apiv1.GET("/articles/:id", article.GetArticle)            // 获取文章详情
		apiv1.GET("/articles", article.ListArticle)               // 获取文章列表
	}
	return r
}
