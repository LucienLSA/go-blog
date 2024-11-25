package main

import (
	"blog-service/conf"
	"blog-service/global"
	"blog-service/models"
	"blog-service/pkg/file"
	"blog-service/pkg/logging"
	"blog-service/routers"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	// 获取当前工作目录
	file.Getwd()
	loading()
	gin.SetMode(global.ServerSetting.RunMode)
	endless.DefaultReadTimeOut = global.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = global.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%v", global.ServerSetting.HttpPort)
	server := endless.NewServer(endPoint, routers.NewRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("server err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Waiting for the system signal to gracefully shutdown.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")

	// s := &http.Server{
	// 	Addr:           ":" + global.ServerSetting.HttpPort, // 根据配置文件
	// 	Handler:        router,
	// 	ReadTimeout:    global.ServerSetting.ReadTimeout,
	// 	WriteTimeout:   global.ServerSetting.WriteTimeout,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.ListenAndServe()
}

// loading 初始化配置、日志和数据库
func loading() {
	// 初始化配置
	err := InitSeting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	// 初始化日志
	err = logging.InitLog()
	if err != nil {
		log.Fatalf("init.setupLog err: %v", err)
	}
	// 初始化数据库
	err = setupMySQL()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

// 初始化配置文件
func InitSeting() error {
	setConf, err := conf.InitConfig()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	// 读取服务端部分
	err = setConf.ReadSection("server", &global.ServerSetting)
	if err != nil {
		return err
	}
	// 读取应用部分
	err = setConf.ReadSection("app", &global.AppSetting)
	if err != nil {
		return err
	}
	// 读取mysql数据库部分
	err = setConf.ReadSection("mysql", &global.MysqlSetting)
	if err != nil {
		return err
	}
	// 读取redis数据库部分
	err = setConf.ReadSection("redis", &global.RedisSetting)
	if err != nil {
		return err
	}
	// 写入和读取超时时间
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

// 初始化MYSQL
func setupMySQL() error {
	var err error
	global.DBEngine, err = models.InitMySQLEngine(global.MysqlSetting)
	if err != nil {
		return err
	}
	return nil
}
