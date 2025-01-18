package controller

import (
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/response"
)

// 放接口文档用到的model
// 接口文档返回的数据格式是一致的，但具体的data类型不一致

type _ResponseUserLogin struct {
	Code    response.ResCode  `json:"code"`    // 业务相应码
	Message string            `json:"message"` // 提示信息
	Data    *models.UserLogin `json:"data"`    //数据
}

type _ResponsePostList struct {
	Code    response.ResCode        `json:"code"`    // 业务相应码
	Message string                  `json:"message"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    //数据
}

type _ResponseCommunityList struct {
	Code    response.ResCode    `json:"code"`    // 业务相应码
	Message string              `json:"message"` // 提示信息
	Data    []*models.Community `json:"data"`    //数据
}

type _ResponseCommunityDetailList struct {
	Code    response.ResCode          `json:"code"`    // 业务相应码
	Message string                    `json:"message"` // 提示信息
	Data    []*models.CommunityDetail `json:"data"`    //数据
}
