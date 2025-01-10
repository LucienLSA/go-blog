package controller

import (
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/request"
	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 用户进行投票
func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			response.ResponseError(c, response.CodeInvalidParam)
			return
		}
		// 翻译并去除掉错误提示中的结构体
		errData := translator.RemoveTopStruct(errs.Translate(translator.Trans))
		response.ResponseErrorMsg(c, response.CodeInvalidParam, errData)
		return
	}
	// 获取当前请求的用户id
	userID, err := request.GetLoginUserID(c)
	if err != nil {
		zap.L().Error("request GetLoginUserID failed", zap.Error(err))
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}
	// 投票业务
	// 将其转化格式，userID和PostID int64转化为string, Kind int8转化为float64
	if err := service.PostVote(userID, p); err != nil {
		zap.L().Error("service.PostVote failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, nil)
}
