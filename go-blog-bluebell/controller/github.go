package controller

import (
	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetGitHubTrendingHandler 获取GitHub热点数据
// @Summary 获取GitHub热点数据
// @Description 获取GitHub热点仓库数据，支持按语言和时间范围筛选
// @Tags GitHub接口
// @Accept json
// @Produce json
// @Param language query string false "编程语言筛选"
// @Param since query string false "时间范围：daily, weekly, monthly"
// @Success 200 {object} types.GitHubTrendingResponse
// @Failure 400 {object} _ResponseError
// @Router /github/trending [get]
func GetGitHubTrendingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取查询参数
		language := c.Query("language")
		since := c.Query("since")

		// 验证时间范围参数
		if since != "" && since != "daily" && since != "weekly" && since != "monthly" {
			e.ResponseErrorMsg(c, e.CodeInvalidParam, "since参数必须是daily、weekly或monthly")
			return
		}

		// 获取GitHub服务
		githubSrv := service.GetGitHubService()

		// 获取热点数据
		result, err := githubSrv.GetGitHubTrendingData(c.Request.Context(), language, since)
		if err != nil {
			zap.L().Error("get github trending data failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}

		e.ResponseSuccessData(c, result)
	}
}

// GetAllGitHubTrendingHandler 获取所有GitHub热点数据
// @Summary 获取所有GitHub热点数据
// @Description 获取所有语言和时间范围的GitHub热点数据
// @Tags GitHub接口
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]types.GitHubTrendingData
// @Failure 400 {object} _ResponseError
// @Router /github/trending/all [get]
func GetAllGitHubTrendingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取GitHub服务
		githubSrv := service.GetGitHubService()

		// 获取所有热点数据
		result, err := githubSrv.GetAllGitHubTrendingData(c.Request.Context())
		if err != nil {
			zap.L().Error("get all github trending data failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}

		e.ResponseSuccessData(c, result)
	}
}

// RefreshGitHubTrendingHandler 手动刷新GitHub热点数据
// @Summary 手动刷新GitHub热点数据
// @Description 手动触发GitHub热点数据刷新任务
// @Tags GitHub接口
// @Accept json
// @Produce json
// @Success 200 {object} _ResponseSuccess
// @Failure 400 {object} _ResponseError
// @Router /github/trending/refresh [post]
func RefreshGitHubTrendingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取GitHub服务
		githubSrv := service.GetGitHubService()

		// 异步执行刷新任务
		go func() {
			githubSrv.RefreshGitHubTrendingData()
		}()

		e.ResponseSuccessData(c, gin.H{
			"message": "GitHub热点数据刷新任务已启动",
		})
	}
}
