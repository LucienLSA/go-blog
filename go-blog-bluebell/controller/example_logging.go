package controller

import (
	"net/http"

	logging "github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExampleLoggingHandler 展示增强日志功能的示例
func ExampleLoggingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userInfo, exists := c.Get("user_info")

		// 获取 TraceID（优先从响应头获取，因为中间件已经设置了）
		traceID := c.Writer.Header().Get("X-Trace-ID")
		if traceID == "" {
			// 如果响应头中没有，尝试从请求头获取
			traceID = c.GetHeader("X-Trace-ID")
		}

		if !exists {
			c.JSON(http.StatusOK, gin.H{
				"message":  "未登录用户",
				"trace_id": traceID,
			})
			return
		}

		// 使用增强的日志记录用户操作
		if user, ok := userInfo.(*ctl.UserInfo); ok {
			// 记录用户操作
			logging.LogUserAction(c.Request.Context(), "example_action", map[string]interface{}{
				"operation": "访问示例接口",
				"user_id":   user.UserId,
				"username":  user.UserName,
				"ip":        c.ClientIP(),
			})

			// 使用带上下文的logger
			ctxLogger := logging.GetLoggerWithContext(c.Request.Context())
			ctxLogger.Info("用户访问示例接口",
				zap.Int64("user_id", user.UserId),
				zap.String("username", user.UserName),
				zap.String("ip", c.ClientIP()),
			)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "已登录用户访问示例接口",
			"trace_id":  traceID,
			"user_info": userInfo,
		})
	}
}
