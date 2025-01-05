package routers

import (
	"net/http"
	"time"

	"github.com/LucienLSA/go-blog/controller"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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
