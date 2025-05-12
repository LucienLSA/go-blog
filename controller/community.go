package controller

import (
	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 查询社区信息
// @Summary 查询社区信息接口
// @Description 查询社区信息返回列表
// @Tags 社区接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} _ResponseCommunityList
// @Router /show/community [get]
// 社区请求
func CommunityListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CommuntityListReq
		// 查询到社区信息（community_id, community_name）以列表形式返回
		if err := ctx.ShouldBind(&req); err != nil {
			zap.L().Error("CommunityList with invalid params", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		l := service.GetCommunitySrv()
		dataList, err := l.GetCommunityList(ctx.Request.Context(), &req)
		if err != nil {
			zap.L().Error("service GetCommunityList failed", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(ctx, dataList)
	}
}

// CommunityHandler 查询社区分类详情信息
// @Summary 查询社区分类详情信息接口
// @Description 查询社区分类详情信息返回列表
// @Tags 社区接口
// @Accept application/json
// @Produce application/json
// @Param community_id path models.ParamCommunityId true "社区ID"
// @Success 200 {object} _ResponseCommunityDetailList
// @Router /community/show/{community_id} [get]
// 社区分类详细请求
func CommunityDetailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CommunityIdReq
		// 1. 获取社区ID
		// cidStr := c.Param("community_id")
		if err := c.ShouldBind(&req); err != nil {
			zap.L().Error("CommunityId with invalid params", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 获取URL参数 并将其字符串参数转化为int64类型
		// communityID, err := strconv.ParseInt(cid, 10, 64)
		// if err != nil {
		// 	zap.L().Error("get community detail with invalid param", zap.Error(err))
		// 	response.ResponseError(c, response.CodeInvalidParam)
		// 	return
		// }
		// 2. 调用服务层获取详情
		l := service.GetCommunitySrv()
		dataList, err := l.GetCommunityDetailList(c.Request.Context(), &req)
		// 3. 返回错误和数据
		if err != nil {
			zap.L().Error("service GetCommunityDetailList failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, dataList)
	}
}

// CommunityHandler 创建社区信息
// @Summary 创建社区信息接口
// @Description 创建社区信息
// @Tags 社区接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object}
// @Router /community/create [post]
// 创建社区
func CreateCommunityHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CommunityCreateResp
		// 查询到社区信息（community_id, community_name）以列表形式返回
		if err := ctx.ShouldBind(&req); err != nil {
			zap.L().Error("CommunityCreate with invalid params", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		l := service.GetCommunitySrv()
		err := l.CreateCommunity(ctx.Request.Context(), &req)
		if err != nil {
			zap.L().Error("service CreateCommunity failed", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(ctx, nil)
	}
}
