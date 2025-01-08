package controller

import (
	"strconv"

	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区请求
func CommunityHandler(c *gin.Context) {
	// 查询到社区信息（community_id, community_name）以列表形式返回
	dataList, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("service GetCommunityList failed", zap.Error(err))
		// 不轻易将服务端报错暴露给外面
		response.ResponseError(c, response.CodeServerBusy)
	}
	response.ResponseSuccessData(c, dataList)
}

// 社区分类详细请求
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区ID
	idStr := c.Param("id")
	communityID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	dataList, err := service.GetCommunityDetailList(communityID)
	if err != nil {
		zap.L().Error("service GetCommunityDetailList failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
	}
	response.ResponseSuccessData(c, dataList)
}
