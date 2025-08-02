package upload

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"go.uber.org/zap"
)

// ValidateImageContent 验证文件内容是否为有效图片
func ValidateImageContent(file multipart.File) error {
	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// 重置文件指针，以便后续使用
	file.Seek(0, 0)

	// 尝试解码图片
	_, _, err = image.Decode(bytes.NewReader(content))
	if err != nil {
		return err
	}

	return nil
}

// 上传到本地文件中
func UploadAvatarToLocalStatic(file multipart.File, userId int64, userName string, fileName string) (filePath string, err error) {
	// 验证图片内容
	if err := ValidateImageContent(file); err != nil {
		zap.L().Error("Invalid image content", zap.Error(err))
		return "", err
	}

	bId := strconv.Itoa(int(userId)) // 路径拼接
	settingPath := settings.Conf.AppConfig.PhotoPathConfig.AvatarPath
	basePath := "." + settingPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}

	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 使用时间戳+用户ID作为文件名，避免重名
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	newFileName := "avatar_" + timestamp + "_" + bId + ext
	avatarPath := basePath + newFileName

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
	return "user" + bId + "/" + newFileName, nil
}

// DeleteOldAvatar 删除旧头像文件
func DeleteOldAvatar(avatarPath string) {
	if avatarPath == "" {
		return
	}

	// 移除开头的斜杠（如果有）
	path := strings.TrimPrefix(avatarPath, "/")

	// 构建完整的文件路径
	settingPath := settings.Conf.AppConfig.PhotoPathConfig.AvatarPath
	fullPath := "." + settingPath + path

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); err == nil {
		// 文件存在，删除它
		err = os.Remove(fullPath)
		if err != nil {
			zap.L().Error("Failed to delete old avatar", zap.String("path", fullPath), zap.Error(err))
		} else {
			zap.L().Info("Successfully deleted old avatar", zap.String("path", fullPath))
		}
	}
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
	return err == nil
}
