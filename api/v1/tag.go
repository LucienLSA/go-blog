package v1

import (
	"blog-service/global"
	"blog-service/pkg/app"
	"blog-service/pkg/e"
	"blog-service/pkg/export"
	"blog-service/pkg/logging"
	"blog-service/pkg/util"
	"blog-service/service/tag_service"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

// type Tag struct{}

// func NewTag() *Tag {
// 	return &Tag{}
// }

// @Summary获取所有标签列表
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
	}
	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: global.AppSetting.DefaultPageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	count, err := tagService.Count()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAGS_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
	// name := c.Query("name")
	// maps := make(map[string]interface{})
	// data := make(map[string]interface{})
	// if name != "" {
	// 	maps["name"] = name
	// }
	// // 获取状态
	// var state int = -1
	// var err error
	// if arg := c.Query("state"); arg != "" {
	// 	state, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		logging.LogrusObj.Infoln(err) // 补充错误处理
	// 	}
	// 	maps["state"] = state
	// }
	// code := e.SUCCESS
	// // 获取标签列表
	// data["lists"], err = models.GetTags(util.GetPage(c), global.AppSetting.DefaultPageSize, maps)
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err) // 补充错误处理
	// 	// return
	// }
	// // 获取标签总数
	// datas, err := models.GetTagTotal(maps)
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err) // 补充错误处理
	// 	// return
	// }
	// data["total"] = datas

	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"data": data,
	// 	"msg":  e.GetMsg(code),
	// })
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query string false "CreatedBy"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"新增错误"}"
// @Router /api/v1/tags [post]
func CreateTags(c *gin.Context) {
	var appG = app.Gin{C: c}
	var form AddTagForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}
	err = tagService.Create()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	// name := c.Query("name")
	// state, err := strconv.Atoi(c.DefaultQuery("state", "0"))
	// if err != nil {
	// name := c.Query("name")
	// state, err := strconv.Atoi(c.DefaultQuery("state", "0"))
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err) // 补充错误处理
	// }
	// createdBy := c.Query("created_by")
	// valid := validation.Validation{}
	// valid.Required(name, "name").Message("名称不能为空")
	// valid.MaxSize(name, 50, "name").Message("名称长度不能超过50")
	// valid.Required(createdBy, "created_by").Message("创建人不能为空")
	// valid.MaxSize(createdBy, 50, "created_by").Message("创建人长度不能超过50")
	// valid.Range(state, 0, 1, "state").Message("状态只能为0或1")

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	if !models.ExitTagByName(name) {
	// 		err = models.CreateTags(name, state, createdBy)
	// 		if err == nil {
	// 			code = e.SUCCESS
	// 		} else {
	// 			code = e.ERROR_ADD_TAG_FAIL
	// 			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 		}
	// 	} else {
	// 		code = e.ERROR_EXIST_TAG
	// 		logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 	}
	// }
	// for _, err := range valid.Errors {
	// 	logging.LogrusObj.Infoln(err)
	// 	// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": make(map[string]string),
	// })

}

type UpdateTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 更新文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {object} body app.Response
// @Failure 500 {object} body app.Response
// @Router /api/v1/tags/{id} [put]
func UpdateTags(c *gin.Context) {
	var appG = app.Gin{C: c}
	Id_1, _ := strconv.Atoi(c.Param("id"))
	var form = UpdateTagForm{ID: Id_1}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = tagService.Update()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err)
	// }
	// name := c.Query("name")
	// modifiedBy := c.Query("modified_by")
	// valid := validation.Validation{}
	// var state int = -1
	// if arg := c.Query("state"); arg != "" {
	// 	state, err = strconv.Atoi(arg)
	// 	if err != nil {
	// 		logging.LogrusObj.Infoln(err)
	// 	}
	// 	valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	// }
	// valid.Required(id, "id").Message("ID不能为空")
	// valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	// valid.Required(name, "name").Message("名称不能为空")
	// valid.MaxSize(modifiedBy, 50, "modified_by").Message("修改人长度不能超过50字符")
	// valid.MaxSize(name, 50, "name").Message("名称长度不能超过50字符")

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	code = e.SUCCESS
	// 	if models.ExistTagByID(id) {
	// 		data := make(map[string]interface{})
	// 		data["modified_by"] = modifiedBy
	// 		if name != "" {
	// 			data["name"] = name
	// 		}
	// 		if state != -1 {
	// 			data["state"] = state
	// 		}
	// 		models.UpdateTags(id, data)
	// 	} else {
	// 		code = e.ERROR_NOT_EXIST_TAG
	// 		logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 	}

	// } else {
	// 	for _, err := range valid.Errors {
	// 		logging.LogrusObj.Infoln(err)
	// 		// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": make(map[string]string),
	// })
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"删除错误"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTags(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id, err := strconv.Atoi(c.Param("id"))
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	if !valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exists {
		logging.LogrusObj.Infoln(err)
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if err := tagService.Delete(); err != nil {
		logging.LogrusObj.Infoln(err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	return
	// }
	// valid := validation.Validation{}
	// valid.Min(id, 1, "id").Message("ID必须大于0")
	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	code = e.SUCCESS
	// 	if models.ExistTagByID(id) {
	// 		models.DeleteTags(id)
	// 	} else {
	// 		code = e.ERROR_NOT_EXIST_TAG
	// 		logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": make(map[string]string),
	// })
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	// state := -1
	if arg := c.PostForm("state"); arg != "" {
		state, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err) // 补充错误处理
		}
		tagService := tag_service.Tag{
			Name:  name,
			State: state,
		}
		filename, err := tagService.Export()
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
			return
		}

		appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"export_url":      export.GetExcelFullUrl(filename),
			"export_save_url": export.GetExcelPath() + filename,
		})
	}
}
