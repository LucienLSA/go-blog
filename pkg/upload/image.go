package upload

import (
	"blog-service/global"
	"blog-service/pkg/file"
	"blog-service/pkg/logging"
	"blog-service/pkg/util"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImageFullUrl(name string) string {
	return global.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return global.AppSetting.ImgSavePath
}

func GetImageFullPath() string {
	return global.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range global.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			fmt.Println("allowExt:", allowExt, "ext:", ext)
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.LogrusObj.Warn(err)
		return false
	}
	return size <= global.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Get error: %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir error: %v", err)
	}
	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission error: %v", err)
	}
	return nil
}
