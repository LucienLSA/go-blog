package controller

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GitHubLoginHandler 重定向到 GitHub 授权页
func GitHubLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		oauth := service.GetOAuthGitHubService()
		// 生成state，并存储到 Redis 10分钟过期
		state, err := oauth.GenerateState(c.Request.Context(), 10*time.Minute)
		if err != nil {
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 构建授权URL，并重定向到 GitHub 授权页
		url := oauth.BuildAuthorizeURL(state)
		c.Redirect(http.StatusFound, url)
	}
}

// GitHubCallbackHandler 处理 GitHub 回调
func GitHubCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			e.ResponseErrorMsg(c, e.CodeInvalidParam, "code/state 缺失")
			return
		}
		oauth := service.GetOAuthGitHubService()
		token, redirectURL, err := oauth.HandleCallback(c.Request.Context(), code, state)
		if err != nil {
			zap.L().Error("oauth github callback failed", zap.Error(err))
			if redirectURL != "" {
				c.Redirect(http.StatusFound, redirectURL)
				return
			}
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		c.Redirect(http.StatusFound, redirectURL+"?token="+token)
	}
}
