package v1

import (
	"blog-service/global"
	"blog-service/models"
	"blog-service/pkg/e"
	"blog-service/pkg/util"
	"log"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

type Tag struct{}

func NewTag() *Tag {
	return &Tag{}
}

// 获取标签列表的具体逻辑
func (t *Tag) GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}
	// 获取状态
	var state int = -1
	var err error
	if arg := c.Query("state"); arg != "" {
		state, err = strconv.Atoi(arg)
		if err != nil {
			maps["state"] = state
		}
	}
	code := e.SUCCESS
	// 获取标签列表
	data["lists"], err = models.GetTags(util.GetPage(c), global.AppSetting.DefaultPageSize, maps)
	if err != nil {
		return
	}
	// 获取标签总数
	datas, err := models.GetTagTotal(maps)
	if err != nil {
		return
	}
	data["total"] = datas

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  e.GetMsg(code),
	})
}

// 新增文章标签
func (t *Tag) CreateTags(c *gin.Context) {
	name := c.Query("name")
	state, err := strconv.Atoi(c.DefaultQuery("state", "0"))
	if err != nil {
		return
	}
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 50, "name").Message("名称长度不能超过50")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 50, "created_by").Message("创建人长度不能超过50")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExitTagByName(name) {
			if models.CreateTags(name, state, createdBy) {
				code = e.SUCCESS
			} else {
				code = e.ERROR_ADD_TAG_FAIL
			}
		} else {
			code = e.ERROR_EXIST_TAG
		}

	} else {
		for _, err := range valid.Errors {
			log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 更新标签
func (t *Tag) UpdateTags(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, err = strconv.Atoi(arg)
		if err != nil {
			return
		}
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 50, "modified_by").Message("修改人长度不能超过50字符")
	valid.MaxSize(name, 50, "name").Message("名称长度不能超过50字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.UpdateTags(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

// 删除标签
func (t *Tag) DeleteTags(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTags(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
