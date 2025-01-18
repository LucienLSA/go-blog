package controller

import (
	"strconv"

	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/service"
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
// @Router /community [get]
// 社区请求
func CommunityHandler(c *gin.Context) {
	// 查询到社区信息（community_id, community_name）以列表形式返回
	dataList, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("service GetCommunityList failed", zap.Error(err))
		// 不轻易将服务端报错暴露给外面
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, dataList)
}

// CommunityHandler 查询社区分类详情信息
// @Summary 查询社区分类详情信息接口
// @Description 查询社区分类详情信息返回列表
// @Tags 社区接口
// @Accept application/json
// @Produce application/json
// @Param community_id path models.ParamCommunityId true "社区ID"
// @Success 200 {object} _ResponseCommunityDetailList
// @Router /community/{community_id} [get]
// 社区分类详细请求
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区ID
	cidStr := c.Param("community_id")
	// 获取URL参数 并将其字符串参数转化为int64类型
	communityID, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		zap.L().Error("get community detail with invalid param", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	// 2. 调用服务层获取详情
	dataList, err := service.GetCommunityDetailList(communityID)
	// 3. 返回错误和数据
	if err != nil {
		zap.L().Error("service GetCommunityDetailList failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, dataList)
}
