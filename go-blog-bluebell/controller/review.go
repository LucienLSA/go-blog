package controller

import (
	"strconv"

	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SubmitPostForReviewHandler 提交帖子进行审核
// @Summary 提交帖子进行审核
// @Description 用户发布帖子，进入审核流程
// @Tags 帖子审核
// @Accept json
// @Produce json
// @Param data body types.PostReviewRequest true "帖子内容"
// @Success 200 {object} types.PostReviewResponse
// @Failure 400 {object} _ResponseError
// @Router /posts/review [post]
func SubmitPostForReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PostReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			e.ResponseErrorMsg(c, e.CodeInvalidParam, "参数错误")
			return
		}
		// 提交帖子进行审核
		resp, err := service.GetReviewSrv().SubmitPostForReview(c.Request.Context(), &req)
		if err != nil {
			zap.L().Error("SubmitPostForReview failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, resp)
	}
}

// GetReviewStatusHandler 查询帖子审核状态
// @Summary 查询帖子审核状态
// @Description 查询帖子审核状态
// @Tags 帖子审核
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Success 200 {object} types.ReviewStatusResponse
// @Failure 400 {object} _ResponseError
// @Router /posts/review/{post_id} [get]
func GetReviewStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("post_id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			e.ResponseErrorMsg(c, e.CodeInvalidParam, "post_id参数错误")
			return
		}
		//  获取帖子审核状态
		resp, err := service.GetReviewSrv().GetReviewStatus(c.Request.Context(), postID)
		if err != nil {
			zap.L().Error("GetReviewStatus failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, resp)
	}
}

// N8nReviewCallbackHandler n8n回调接口
// @Summary n8n回调接口
// @Description n8n审核完成后回调此接口
// @Tags 帖子审核
// @Accept json
// @Produce json
// @Param data body types.ReviewCallbackRequest true "n8n回调数据"
// @Success 200 {object} types.ReviewCallbackResponse
// @Failure 400 {object} _ResponseError
// @Router /posts/review/callback [post]
func N8nReviewCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ReviewCallbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			e.ResponseErrorMsg(c, e.CodeInvalidParam, "参数错误")
			return
		}
		// 处理n8n回调
		resp, err := service.GetReviewSrv().ProcessN8nCallback(c.Request.Context(), &req)
		if err != nil {
			zap.L().Error("N8nReviewCallback failed", zap.Error(err))
			e.ResponseErrorMsg(c, e.CodeServerBusy, "审核回调处理失败")
			return
		}
		e.ResponseSuccessData(c, resp)
	}
}

// GetReviewStatisticsHandler 获取审核统计
func GetReviewStatisticsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := service.GetReviewSrv().GetReviewStatistics(c.Request.Context())
		if err != nil {
			zap.L().Error("GetReviewStatistics failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, resp)
	}
}

// GetPostsByReviewStatusHandler 根据审核状态获取帖子列表
func GetPostsByReviewStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		pageNumStr := c.DefaultQuery("page_num", "1")
		pageSizeStr := c.DefaultQuery("page_size", "10")
		pageNum, _ := strconv.ParseInt(pageNumStr, 10, 64)
		pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
		resp, err := service.GetReviewSrv().GetPostsByReviewStatus(c.Request.Context(), status, pageNum, pageSize)
		if err != nil {
			zap.L().Error("GetPostsByReviewStatus failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, resp)
	}
}
