package upload

import (
	"mime/multipart"
	"strconv"
)

func UploadAvatarToLocalStatic(file multipart.File, userId int64, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(userId)) // 路径拼接
}
