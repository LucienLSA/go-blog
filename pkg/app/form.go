package app

import (
	"blog-service/pkg/e"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	return http.StatusOK, e.SUCCESS
}
