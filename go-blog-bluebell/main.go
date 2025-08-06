package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	logging "github.com/LucienLSA/go-blog/logger"
// 	"github.com/LucienLSA/go-blog/pkg/scheduler"
// 	"github.com/LucienLSA/go-blog/pkg/snowflake"
// 	"github.com/LucienLSA/go-blog/pkg/translator"
// 	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
// 	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
// 	"github.com/LucienLSA/go-blog/routers"
// 	"github.com/LucienLSA/go-blog/settings"
// 	"github.com/spf13/viper"
// 	"go.uber.org/zap"
// )

// // @title bluebell go-web
// // @version 1.0
// // @description This is a go-web project from https://www.bilibili.com/cheese/play/ep265306
// // @termsOfService http://swagger.io/terms/

// // @contact.name lucien
// // @contact.url none
// // @contact.email none

// // @license.name Apache 2.0
// // @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// // @host 127.0.0.1:8081
// // @BasePath /api/v1/

// func main() {
// 	// 1. 加载配置
// 	if err := settings.InitSettings(); err != nil {
// 		fmt.Printf("init settings failed, err:%v\n", err)
// 		return
// 	}
// 	sConf := settings.Conf
// 	zap.L().Info("init settings success")

// 	// 2. 初始化日志
// 	if err := logging.InitLogger(sConf.LogConfig, sConf.Mode); err != nil {
// 		fmt.Printf("init logger failed, err:%v\n", err)
// 		return
// 	}
// 	zap.L().Info("init logger success")
// 	// 延迟将缓存区的日志追加
// 	defer zap.L().Sync()

// 	// 3. 初始化MySQL
// 	if err := mysql.InitMysql(sConf.MySQLConfig); err != nil {
// 		fmt.Printf("init mysql failed, err:%v\n", err)
// 		return
// 	}
// 	zap.L().Info("init mysql success")

// 	// 4. 初始化Redis链接
// 	if err := redisCache.InitRedis(sConf.RedisConfig); err != nil {
// 		fmt.Printf("init redis failed, err:%v\n", err)
// 		return
// 	}
// 	zap.L().Info("init redis success")
// 	defer redisCache.Close()

// 	// 5. 初始化雪花算法
// 	if err := snowflake.InitSnowflake(sConf.AppConfig.StartTime, sConf.AppConfig.MachineID); err != nil {
// 		fmt.Printf("init snowflake failed, err:%v\n", err)
// 		return
// 	}
// 	zap.L().Info("init snowflake success")

// 	// 6. 注册路由
// 	// 初始化gin框架内置的校验器 翻译校验错误信息
// 	if err := translator.InitTrans("zh"); err != nil {
// 		fmt.Printf("init validator translator failed, err:%v\n", err)
// 		return
// 	}
// 	r := routers.SetupRouter(sConf.Mode)

// 	// 7. 启动GitHub热点数据定时任务
// 	scheduler.StartGitHubTrendingScheduler()
// 	zap.L().Info("github trending data scheduler started")

// 	// 8. 启动服务（优雅关机）
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
// 		Handler: r,
// 	}

// 	go func() {
// 		// 开启一个goroutine启动服务
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			zap.L().Fatal("listen: %s\n", zap.Error(err))
// 		}
// 	}()

// 	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
// 	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
// 	// kill 默认会发送 syscall.SIGTERM 信号
// 	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
// 	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
// 	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
// 	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
// 	zap.L().Info("Shutdown Server ...")

// 	// 停止GitHub热点数据定时任务
// 	scheduler.StopGitHubTrendingScheduler()
// 	zap.L().Info("github trending data scheduler stopped")

// 	// 创建一个5秒超时的context
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
// 	if err := srv.Shutdown(ctx); err != nil {
// 		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
// 	}
// 	zap.L().Info("Server exiting")
// }
