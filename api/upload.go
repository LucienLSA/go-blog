package api

import (
	"blog-service/pkg/app"
	"blog-service/pkg/e"
	"blog-service/pkg/logging"
	"blog-service/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 上传图片
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.LogrusObj.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if image == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}
	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.LogrusObj.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}
	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.LogrusObj.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_url":       upload.GetImageFullUrl(imageName),
		"image_save_path": savePath + imageName,
	})

	// code := e.SUCCESS
	// data := make(map[string]string)

	// file, image, err := c.Request.FormFile("image")
	// if err != nil {
	// 	logging.LogrusObj.Warn(err)
	// 	code = e.ERROR
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": code,
	// 		"msg":  e.GetMsg(code),
	// 		"data": data,
	// 	})
	// }
	// if image == nil {
	// 	code = e.INVALID_PARAMS
	// 	logging.LogrusObj.Info(e.GetMsg(code))
	// } else {
	// 	imageName := upload.GetImageName(image.Filename)
	// 	fullpath := upload.GetImageFullPath()
	// 	savePath := upload.GetImagePath()

	// 	src := fullpath + imageName
	// 	fmt.Println(src)
	// 	if upload.CheckImageExt(imageName) && upload.CheckImageSize(file) {
	// 		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
	// 		logging.LogrusObj.Info(e.GetMsg(code))
	// 	} else {
	// 		err := upload.CheckImage(fullpath)
	// 		if err != nil {
	// 			logging.LogrusObj.Warn(err)
	// 			code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
	// 			logging.LogrusObj.Info(e.GetMsg(code))
	// 		} else if err := c.SaveUploadedFile(image, src); err != nil {
	// 			logging.LogrusObj.Warn(err)
	// 			code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
	// 			logging.LogrusObj.Info(e.GetMsg(code))
	// 		} else {
	// 			data["image_url"] = upload.GetImageFullUrl(imageName)
	// 			fmt.Println(data["image_url"])
	// 			data["image_save_path"] = savePath + imageName
	// 			fmt.Println(data["image_save_path"])
	// 		}
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": data,
	// })
}
