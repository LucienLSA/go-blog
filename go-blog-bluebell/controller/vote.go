package controller

import (
	"github.com/LucienLSA/go-blog/pkg/e"

	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteHandler 用户投票
// @Summary 用户投票
// @Description 用户对帖子进行投票（点赞/点踩）
// @Tags 帖子接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body types.PostVoteDataReq true "投票参数"
// @Success 200 {object} _ResponseSuccess
// @Failure 400 {object} _ResponseError
// @Router /posts/vote [post]
func PostVoteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 参数校验
		// p := new(models.ParamVoteData)
		var req types.PostVoteDataReq
		if err := c.ShouldBindJSON(&req); err != nil {
			errs, ok := err.(validator.ValidationErrors) // 类型断言
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			// 翻译并去除掉错误提示中的结构体
			errData := translator.RemoveTopStruct(errs.Translate(translator.Trans))
			e.ResponseErrorMsg(c, e.CodeInvalidParam, errData)
			return
		}
		// 获取当前请求的用户id

		// userID, err := request.GetLoginUserID(c)
		// if err != nil {
		// 	zap.L().Error("request GetLoginUserID failed", zap.Error(err))
		// 	e.ResponseError(c, e.CodeNeedLogin)
		// 	return
		// }
		// 投票业务
		// 将其转化格式，userID和PostID int64转化为string, Kind int8转化为float64

		l := service.GetVoteSrv()
		if err := l.PostVote(c.Request.Context(), &req); err != nil {
			zap.L().Error("service.PostVote failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, nil)
	}
}
