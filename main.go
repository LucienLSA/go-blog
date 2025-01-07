package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/dao/redis"
	"github.com/LucienLSA/go-blog/logger"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/routers"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// go web 开发通用的脚手架模板

func main() {
	// 1. 加载配置
	if err := settings.InitSettings(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	zap.L().Info("init settings success")

	// 2. 初始化日志
	if err := logger.InitLogger(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	zap.L().Info("init logger success")
	// 延迟将缓存区的日志追加
	defer zap.L().Sync()

	// 3. 初始化MySQL
	if err := mysql.InitMysql(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	zap.L().Info("init mysql success")
	defer mysql.Close()

	// 4. 初始化Redis链接
	if err := redis.InitRedis(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	zap.L().Info("init redis success")
	defer redis.Close()

	// 5. 初始化雪花算法
	if err := snowflake.InitSnowflake(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	zap.L().Info("init snowflake success")

	// 6. 注册路由
	// 初始化gin框架内置的校验器 翻译校验错误信息
	if err := translator.InitTrans("zh"); err != nil {
		fmt.Printf("init validator translator failed, err:%v\n", err)
		return
	}
	r := routers.SetupRouter(settings.Conf.Mode)

	// 7. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}

func startSever() {

}
