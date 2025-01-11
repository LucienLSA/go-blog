package controller

import (
	"strconv"

	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/request"
	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 发布帖子
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数 参数校验
	p := new(models.Post)
	// 如果请求中有json格式数据，才使用shouldbindjson
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
		zap.L().Error("request GetLoginUserID failed", zap.Error(err))
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

// 新版查询帖子，根据前端传来的参数动态获取帖子列表
// 按照创建时间或者分数排序
// 1. 获取参数
// GET请求参数：/api/v1/searchposts?page_num=1&page_size=5&order=time
// 以Query获取参数
// 2. redis查询id列表
// 3. 根据id去数据库查询帖子详细信息
func SearchPostListHandler(c *gin.Context) {
	// 获取分页参数
	// c.ShouldBind() //根据请求数据类型选择相应的方法获取数据
	// 初始化结构体传入初始参数
	p := models.ParamPostList{
		PageNum:  settings.Conf.AppConfig.PageNum,
		PageSize: settings.Conf.AppConfig.PageSize,
		Order:    models.OrderTime,
	}
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("SearchPostListHandler failed with invalid params", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParam)
		return
	}
	dataList, err := service.SearchPostList(&p)
	if err != nil {
		zap.L().Error("service SearchPostList failed", zap.Error(err))
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccessData(c, dataList)
}
