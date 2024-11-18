package util

import (
	"blog-service/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	result := 0
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err == nil && page > 0 {
		result = (page - 1) * global.AppSetting.DefaultPageSize
	}

	return result
}
