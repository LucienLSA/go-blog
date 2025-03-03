package upload

import (
	"github.com/LucienLSA/go-blog/settings"
)

// AvatarURL 头像地址
func AvatarURL() string {
	if settings.Conf.AppConfig.UploadModel == settings.UploadModelOss {
		return "http://"
	} else {
		pConfig := settings.Conf.AppConfig
		return pConfig.PhotoHost + ":" + settings.Conf.AppConfig.Port + pConfig.AvatarPath
	}
}
