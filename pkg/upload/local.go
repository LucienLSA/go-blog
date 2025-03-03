package upload

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/LucienLSA/go-blog/settings"
	"go.uber.org/zap"
)

// 上传到本地文件中
func UploadAvatarToLocalStatic(file multipart.File, userId int64, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId)) // 路径拼接
	settingPath := settings.Conf.AppConfig.PhotoPathConfig.AvatarPath
	basePath := "." + settingPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userName + ".jpg" // 将file的后缀提取出来
	content, err := ioutil.ReadAll(file)
	if err != nil {
		zap.L().Error("avatar content read failed", zap.Error(err))
		return "", err
	}
	err = ioutil.WriteFile(avatarPath, content, 0666)
	if err != nil {
		zap.L().Error("avatar content write failed", zap.Error(err))
		return "", err
	}
	return "user" + bId + "/" + userName + ".jpg", nil
}

func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 755)
	if err != nil {
		return false
	}
	return true
}
