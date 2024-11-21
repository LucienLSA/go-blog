package v1

import (
	"blog-service/global"
	"blog-service/models"
	"blog-service/pkg/e"
	"blog-service/pkg/logging"
	"blog-service/pkg/util"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

// 获取单个文章
func (a *Article) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.LogrusObj.Infoln(err)
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data, err = models.GetArticle(id)
			if err != nil {
				logging.LogrusObj.Infoln(err)
			}
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
		}
	} else {
		for _, err := range valid.Errors {
			logging.LogrusObj.Infoln(err)
			// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  e.GetMsg(code),
	})

}

// 获取多个文章
func (a *Article) ListArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	// var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err)
		}
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("state只能为0或1")
	}

	// var tagID int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err)
		}
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("tag_id必须大于0")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		datas, err := models.GetArticles(util.GetPage(c), global.AppSetting.DefaultPageSize, maps)
		if err != nil {
			logging.LogrusObj.Infoln(err) // 补充错误处理
		}
		data["articles"] = datas
		totals, err := models.GetArticleTotal(maps)
		if err != nil {
			logging.LogrusObj.Infoln(err) // 补充错误处理
		}
		data["total"] = totals
	} else {
		for _, err := range valid.Errors {
			logging.LogrusObj.Infoln(err)
			// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  e.GetMsg(code),
	})

}

// 新增文章
func (a *Article) CreateArticles(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Query("tag_id"))
	if err != nil {
		logging.LogrusObj.Infoln(err) // 补充错误处理
	}
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state, err := strconv.Atoi(c.DefaultQuery("state", "0"))
	if err != nil {
		logging.LogrusObj.Infoln(err) // 补充错误处理
	}
	valid := validation.Validation{}
	valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("描述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建者不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagID) {
			data := map[string]interface{}{
				"tag_id":     tagID,
				"title":      title,
				"desc":       desc,
				"content":    content,
				"created_by": createdBy,
				"state":      state,
			}
			if models.AddArticle(data) {
				code = e.SUCCESS
			} else {
				code = e.ERROR_ADD_ARTICLE_FAIL
				logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
			}

		} else {
			code = e.ERROR_NOT_EXIST_TAG
			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
		}
	} else {
		for _, err := range valid.Errors {
			logging.LogrusObj.Infoln(err)
			// log.Panicln("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// 更新文章
func (a *Article) UpdateArticles(c *gin.Context) {
	valid := validation.Validation{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.LogrusObj.Infoln(err)
	}
	tagID, err := strconv.Atoi(c.Query("tag_id"))
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	// var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, err := strconv.Atoi(arg)
		if err != nil {
			logging.LogrusObj.Infoln(err)
		}
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}
	valid.Min(id, 1, "id").Message("id必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题长度不能超过100字符")
	valid.MaxSize(desc, 255, "desc").Message("描述长度不能超过200字符")
	valid.MaxSize(content, 65535, "content").Message("内容长度不能超过65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改者不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改者长度不能超过100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagID) {
				data := make(map[string]interface{})
				if tagID > 0 {
					data["tag_id"] = tagID
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}
				data["modified_by"] = modifiedBy
				models.UpdateArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
				logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
			}
		}
	} else {
		for _, err := range valid.Errors {
			logging.LogrusObj.Infoln(err)
			// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// 删除文章
func (a *Article) DeleteArticles(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.LogrusObj.Infoln(err)
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
		}
	} else {
		for _, err := range valid.Errors {
			logging.LogrusObj.Infoln(err)
			// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}
