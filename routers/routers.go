package routers

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/controller"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin设置发布模式
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ping", func(ctx *gin.Context) {
		time.Sleep(5 * time.Second)
		ctx.String(http.StatusOK, "pong ok")
	})

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	return r
}
