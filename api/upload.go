package api

import (
	"blog-service/pkg/e"
	"blog-service/pkg/logging"
	"blog-service/pkg/upload"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.LogrusObj.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}
	if image == nil {
		code = e.INVALID_PARAMS
		logging.LogrusObj.Info(e.GetMsg(code))
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullpath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullpath + imageName
		fmt.Println(src)
		if upload.CheckImageExt(imageName) && upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
			logging.LogrusObj.Info(e.GetMsg(code))
		} else {
			err := upload.CheckImage(fullpath)
			if err != nil {
				logging.LogrusObj.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				logging.LogrusObj.Info(e.GetMsg(code))
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.LogrusObj.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
				logging.LogrusObj.Info(e.GetMsg(code))
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				fmt.Println(data["image_url"])
				data["image_save_path"] = savePath + imageName
				fmt.Println(data["image_save_path"])
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
