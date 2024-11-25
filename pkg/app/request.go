package app

import (
	"blog-service/pkg/logging"

	"github.com/beego/beego/v2/core/validation"
)

// 记录错误 提前返回
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.LogrusObj.Info(err.Key + " " + err.Message)
	}
	return
}
