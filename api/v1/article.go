package v1

import (
	"blog-service/global"
	"blog-service/pkg/app"
	"blog-service/pkg/e"
	"blog-service/pkg/logging"
	"blog-service/pkg/util"
	"blog-service/service/article_service"
	"blog-service/service/tag_service"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

// type Article struct {
// }

// func NewArticle() *Article {
// 	return &Article{}
// }

// @Summary 获取单篇文章详情
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"获取错误"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.LogrusObj.Infoln(err)
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)

	// code := e.INVALID_PARAMS
	// var data interface{}
	// if !valid.HasErrors() {
	// 	if models.ExistArticleByID(id) {
	// 		data, err = models.GetArticle(id)
	// 		if err != nil {
	// 			logging.LogrusObj.Infoln(err)
	// 		}
	// 		code = e.SUCCESS
	// 	} else {
	// 		code = e.ERROR_NOT_EXIST_ARTICLE
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
	// 	"data": data,
	// 	"msg":  e.GetMsg(code),
	// })
}

// @Summary 获取所有文章列表
// @Produce  json
// @Param tag_id query int false "TagID"
// @Param state query string false "State"
// @Param created_by query string false "CreatedBy"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"获取错误"}"
// @Router /api/v1/articles [get]
func ListArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	// data := make(map[string]interface{})
	// maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.PostForm("state"); arg != "" {
		state, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err)
		}
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}
	var tagID int = -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err)
		}
		valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagID,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: global.AppSetting.DefaultPageSize,
	}

	totol, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
	}

	data := make(map[string]interface{})
	data["articles"] = articles
	data["total"] = totol
	appG.Response(http.StatusOK, e.SUCCESS, data)

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	code = e.SUCCESS
	// 	datas, err := models.GetArticles(util.GetPage(c), 10, maps)
	// 	if err != nil {
	// 		logging.LogrusObj.Infoln(err) // 补充错误处理
	// 	}
	// 	data["articles"] = datas
	// 	totals, err := models.GetArticleTotal(maps)
	// 	if err != nil {
	// 		logging.LogrusObj.Infoln(err) // 补充错误处理
	// 	}
	// 	data["total"] = totals
	// } else {
	// 	for _, err := range valid.Errors {
	// 		logging.LogrusObj.Infoln(err)
	// 		// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"data": data,
	// 	"msg":  e.GetMsg(code),
	// })

}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 新增文章
// @Produce  json
// @Param tag_id query int true "TagID"
// @Param title query string true "Title"
// @Param desc query string true "Desc"
// @Param content query string true "Content"
// @Param created_by query string true "CreatedBy"
// @Param state query int true "State"
// @Param cover_image_url query string true "Cover_image_url"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"新增错误"}"
// @Router /api/v1/articles [post]
func CreateArticles(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)

	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CreatedBy:     form.CreatedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err := articleService.Create(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

	// tagID, err := strconv.Atoi(c.Query("tag_id"))
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err) // 补充错误处理
	// }
	// title := c.Query("title")
	// desc := c.Query("desc")
	// content := c.Query("content")
	// createdBy := c.Query("created_by")
	// state, err := strconv.Atoi(c.DefaultQuery("state", "0"))
	// coverImageUrl := c.Query("cover_image_url")
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err) // 补充错误处理
	// }
	// valid := validation.Validation{}
	// valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	// valid.Required(title, "title").Message("标题不能为空")
	// valid.Required(desc, "desc").Message("描述不能为空")
	// valid.Required(content, "content").Message("内容不能为空")
	// valid.Required(createdBy, "created_by").Message("创建者不能为空")
	// valid.Required(coverImageUrl, "cover_image_url").Message("封面图片不能为空")
	// // valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("封面图片长度不能超过255字符")
	// valid.Range(state, 0, 1, "state").Message("状态只能为0或1")

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	if models.ExistTagByID(tagID) {
	// 		data := map[string]interface{}{
	// 			"tag_id":          tagID,
	// 			"title":           title,
	// 			"desc":            desc,
	// 			"content":         content,
	// 			"created_by":      createdBy,
	// 			"state":           state,
	// 			"cover_image_url": coverImageUrl,
	// 		}
	// 		if models.AddArticle(data) {
	// 			code = e.SUCCESS
	// 		} else {
	// 			code = e.ERROR_ADD_ARTICLE_FAIL
	// 			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 		}

	// 	} else {
	// 		code = e.ERROR_NOT_EXIST_TAG
	// 		logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 	}
	// } else {
	// 	for _, err := range valid.Errors {
	// 		logging.LogrusObj.Infoln(err)
	// 		// log.Panicln("err.key: %s, err.message: %s", err.Key, err.Message)
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": make(map[string]interface{}),
	// })
}

type UpdateArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary 更新文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Param cover_image_url query string true "Cover_image_url"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"更新错误"}"package v1
// @Router /api/v1/articles/{id} [put]
func UpdateArticles(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
	}
	var form = UpdateArticleForm{ID: id}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}

	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = articleService.Update()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

	// valid := validation.Validation{}

	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err)
	// }
	// tagID, err := strconv.Atoi(c.Query("tag_id"))
	// title := c.Query("title")
	// desc := c.Query("desc")
	// content := c.Query("content")
	// modifiedBy := c.Query("modified_by")
	// coverImageUrl := c.Query("cover_image_url")
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err)
	// }
	// // var state int = -1
	// if arg := c.Query("state"); arg != "" {
	// 	state, err := strconv.Atoi(arg)
	// 	if err != nil {
	// 		logging.LogrusObj.Infoln(err)
	// 	}
	// 	valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	// }
	// valid.Min(id, 1, "id").Message("id必须大于0")
	// valid.MaxSize(title, 100, "title").Message("标题长度不能超过100字符")
	// valid.MaxSize(desc, 255, "desc").Message("描述长度不能超过255字符")
	// valid.MaxSize(content, 65535, "content").Message("内容长度不能超过65535字符")
	// // valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("封面图片长度不能超过255字符")
	// valid.Required(coverImageUrl, "cover_image_url").Message("封面图片不能为空")
	// valid.Required(modifiedBy, "modified_by").Message("修改者不能为空")
	// valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改者长度不能超过100字符")

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	if models.ExistArticleByID(id) {
	// 		if models.ExistTagByID(tagID) {
	// 			data := make(map[string]interface{})
	// 			if tagID > 0 {
	// 				data["tag_id"] = tagID
	// 			}
	// 			if title != "" {
	// 				data["title"] = title
	// 			}
	// 			if desc != "" {
	// 				data["desc"] = desc
	// 			}
	// 			if content != "" {
	// 				data["content"] = content
	// 			}

	// 			data["modified_by"] = modifiedBy
	// 			models.UpdateArticle(id, data)
	// 			code = e.SUCCESS
	// 		} else {
	// 			code = e.ERROR_NOT_EXIST_ARTICLE
	// 			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 		}
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
	// 	"data": make(map[string]interface{}),
	// })
}

// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":e.SUCCESS,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":e.ERROR,"data":{},"msg":"删除错误"}"package v1
// @Router /api/v1/articles/{id} [delete]
func DeleteArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
	}
	valid.Min(id, 1, "id").Message("id必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	logging.LogrusObj.Infoln(err)
	// }
	// valid := validation.Validation{}
	// valid.Min(id, 1, "id").Message("id必须大于0")

	// code := e.INVALID_PARAMS
	// if !valid.HasErrors() {
	// 	if models.ExistArticleByID(id) {
	// 		models.DeleteArticle(id)
	// 		code = e.SUCCESS
	// 	} else {
	// 		code = e.ERROR_NOT_EXIST_ARTICLE
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
	// 	"data": make(map[string]interface{}),
	// })
}
