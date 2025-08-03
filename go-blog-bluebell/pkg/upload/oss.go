package upload

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"go.uber.org/zap"
)

// UploadToQiNiu 封装上传图片到七牛云然后返回状态和图片的url，单张
func UploadToQiNiuAvatar(file multipart.File, userName string, fileSize int64, fileName string) (path string, err error) {
	// 验证图片内容
	if err := ValidateImageContent(file); err != nil {
		zap.L().Error("Invalid image content for OSS upload", zap.Error(err))
		return "", err
	}

	qConfig := settings.Conf.OssConfig
	var AccessKey = qConfig.AccessKeyId
	var SerectKey = qConfig.AccessKeySecret
	var Bucket = qConfig.BucketName
	var ImgUrl = qConfig.QiNiuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}

	mac := credentials.NewCredentials(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 使用时间戳+用户名作为key，避免重名
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	key := "avatar/" + userName + "_" + timestamp + ext

	err = formUploader.Put(context.Background(), &ret, upToken, key, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + "/" + key
	return url, nil
}
