package controller

import (
	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler 发布帖子
// @Summary 发布帖子
// @Description 发布帖子
// @Tags 帖子接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body models.ParamPostCreate false "请求参数"
// @Security ApiKeyAuth
// @Router /posts/post [post]
// 发布帖子
func CreatePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取参数 参数校验
		// p := new(models.Post)
		var req types.PostCreateReq
		// 如果请求中有json格式数据，才使用shouldbindjson
		if err := c.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("create post with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// 2. 创建帖子
		l := service.GetPostSrv()
		if err := l.CreatePost(c, &req); err != nil {
			zap.L().Error("server CreatePost failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 3. 返回响应
		e.ResponseSuccessData(c, nil)
	}
}

// GetPostListHandler 获取帖子列表分页展示
// @Summary 获取帖子列表分页展示接口
// @Description 根据分页参数获取帖子列表分页展示
// @Tags 帖子接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostSearch false "请求参数"
// @Success 200 {object} _ResponsePostList
// @Router /posts/show [get]
// 获取帖子列表分页展示
func GetPostListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PostSearchReq
		if err := c.ShouldBind(&req); err != nil {
			zap.L().Error("GetPostList with invalid params", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(c, e.CodeInvalidParam)
			return
		}
		// 获取分页参数
		// pageNumStr := req.PageNum
		// pageSizeStr := req.PageSize

		// pageNum, pageSize, err := request.GetPageInfo(c)
		// if err != nil {
		// 	zap.L().Error("request GetPageInfo failed", zap.Error(err))
		// 	response.ResponseError(c, response.CodeServerBusy)
		// 	return
		// }
		// 查询帖子以列表形式返回数据
		l := service.GetPostSrv()
		dataList, err := l.GetPostList(c, req.PageNum, req.PageSize)
		if err != nil {
			zap.L().Error("service GetPostListHandler failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, dataList)
	}
}

// GetPostDetailHandler 获取帖子分类详情
// @Summary 获取帖子分类详情接口
// @Description 根据帖子id获取帖子分类详情接口
// @Tags 帖子接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param post_id path models.ParamPostId true "帖子ID"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts/{post_id} [get]
// 获取帖子分类详情
func GetPostDetailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PostIdReq
		// 1. 获取帖子ID
		// pIdStr := c.Param("post_id")
		// 获取URL参数 并将其字符串参数转化为int64类型
		// postID, err := strconv.ParseInt(pIdStr, 10, 64)
		// if err != nil {
		// 	zap.L().Error("get post detail with invalid param", zap.Error(err))
		// 	e.ResponseError(c, e.CodeInvalidParam)
		// 	return
		// }
		if err := c.ShouldBind(&req); err != nil {
			zap.L().Error("GetPostList with invalid params", zap.Error(err))
			// 不轻易将服务端报错暴露给外面
			e.ResponseError(c, e.CodeInvalidParam)
			return
		}
		// 2. 调用服务层获取帖子详情
		l := service.GetPostSrv()
		dataList, err := l.GetPostDetailList(c, &req)
		// fmt.Println(dataList.Post.CreateTime, dataList.Post.UpdateTime)
		// 3. 返回错误和数据
		if err != nil {
			zap.L().Error("service GetPostDetailList failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, dataList)
	}
}

// SearchPostListHandler 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query types.PostListReq false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts/search [get]
// 新版查询帖子，根据前端传来的参数动态获取帖子列表
// 按照创建时间或者分数排序
// 1. 获取参数
// GET请求参数：/api/v1//posts/search?page_num=1&page_size=5&order=time
// 以Query获取参数
// 2. redis查询id列表
// 3. 根据id去数据库查询帖子详细信息
func SearchPostListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取分页参数
		// c.ShouldBind() //根据请求数据类型选择相应的方法获取数据
		// 初始化结构体传入初始参数
		var req types.PostListReq
		req = types.PostListReq{
			PageNum:  settings.Conf.AppConfig.PageNum,
			PageSize: settings.Conf.AppConfig.PageSize,
			Order:    types.OrderTime,
			// CommunityID不初始化，根据有无查询的业务不一样
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			zap.L().Error("SearchPostListHandler failed with invalid params", zap.Error(err))
			e.ResponseError(c, e.CodeInvalidParam)
			return
		}
		// dataList, err := service.SearchPostList(&p)
		l := service.GetPostSrv()
		dataList, err := l.GetPostListNew(c, &req)
		if err != nil {
			zap.L().Error("service GetPostListNew failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(c, dataList)
	}
}

// // // 根据社区查询帖子信息 整合到一个handler中
// // func CommunityPostListHandler(c *gin.Context) {
// // 	// 与 SearchPostListHandler 类似
// // 	cIdStr := c.Param("community_id")
// // 	// 获取URL参数 并将其字符串参数转化为int64类型
// // 	communityID, err := strconv.ParseInt(cIdStr, 10, 64)
// // 	if err != nil {
// // 		zap.L().Error("get post detail with invalid param", zap.Error(err))
// // 		response.ResponseError(c, response.CodeInvalidParam)
// // 		return
// // 	}
// // 	p := &models.ParamCommunityPostList{
// // 		ParamPostList: &models.ParamPostList{
// // 			PageNum:  settings.Conf.AppConfig.PageNum,
// // 			PageSize: settings.Conf.AppConfig.PageSize,
// // 			Order:    models.OrderTime,
// // 		}, CommunityID: communityID,
// // 	}
// // 	if err := c.ShouldBindQuery(&p); err != nil {
// // 		zap.L().Error("CommunityPostListHandler failed with invalid params", zap.Error(err))
// // 		response.ResponseError(c, response.CodeInvalidParam)
// // 		return
// // 	}
// // 	dataList, err := service.CommunitySearchPostList(p)
// // 	if err != nil {
// // 		zap.L().Error("service CommunityPostList failed", zap.Error(err))
// // 		response.ResponseError(c, response.CodeServerBusy)
// // 		return
// // 	}
// // 	response.ResponseSuccessData(c, dataList)
// // }
