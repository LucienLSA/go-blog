package controller

import (
	"strconv"

	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/request"
	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 发布帖子
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数 参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误 直接返回响应
		zap.L().Error("create post with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.ResponseError(c, response.CodeInvalidParam)
			return
		}
		response.ResponseErrorMsg(c, response.CodeInvalidParam,
			translator.RemoveTopStruct(errs.Translate(translator.Trans)))
		return
	}
	// 从c 取到当前发请求的用户userid
	userID, err := request.GetLoginUserID(c)
	if err != nil {
		response.ResponseError(c, response.CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := service.CreatePost(p); err != nil {
		zap.L().Error("server CreatePost failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.ResponseSuccessData(c, nil)
}

// 获取帖子列表分页展示
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	pageNum, pageSize, err := request.GetPageInfo(c)
	if err != nil {
		zap.L().Error("request GetPageInfo failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	// 查询帖子以列表形式返回数据
	dataList, err := service.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("service GetPostListHandler failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, dataList)
}

// 获取帖子分类详情
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取帖子ID
	pIdStr := c.Param("post_id")
	// 获取URL参数 并将其字符串参数转化为int64类型
	postID, err := strconv.ParseInt(pIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	// 2. 调用服务层获取帖子详情
	dataList, err := service.GetPostDetailList(postID)
	// fmt.Println(dataList.Post.CreateTime, dataList.Post.UpdateTime)
	// 3. 返回错误和数据
	if err != nil {
		zap.L().Error("service GetPostDetailList failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, dataList)
}
