package logging

import (
	"blog-service/pkg/file"
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var LogrusObj *logrus.Logger

func InitLog() error {
	if LogrusObj != nil {
		src, _ := setOutPutFile()
		LogrusObj.Out = src
		return nil
	}
	// 实例化
	logger := logrus.New()
	src, _ := setOutPutFile()
	logger.Out = src                   // 设置输出
	logger.SetLevel(logrus.DebugLevel) // 设置日志级别
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
	return nil
}

func setOutPutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := file.Getwd(); err == nil { // 获取工作目录
		logFilePath = path.Join(dir, "runtime", "logs") + string(os.PathSeparator)
	}
	_, err := os.Stat(logFilePath)
	// if os.IsNotExist(err) {
	// 	if err := os.MkdirAll(logFilePath, 0777); err != nil {
	// 		log.Println(err.Error())
	// 		return nil, err
	// 	}
	// }
	err = file.IsNotExistMkDir(logFilePath)
	if err != nil {
		return nil, err
	}
	logFileName := now.Format("2006-01-02") + ".log"
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	// 写入文件
	// src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	//
	src, err := file.Open(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
