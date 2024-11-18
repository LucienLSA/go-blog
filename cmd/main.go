package main

import (
	"blog-service/global"
	"blog-service/internal/model"
	"blog-service/internal/routers"
	"blog-service/pkg/setting"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化配置
	err := setupConfig()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort, // 根据配置文件
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

// 初始化配置
func setupConfig() error {
	setting, err := setting.NewConfig()
	if err != nil {
		return err
	}
	// 读取服务端部分
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	// 读取应用部分
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	// 读取数据库部分
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	// 写入和读取超时时间
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

func setupLogger() error {

}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.DBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
