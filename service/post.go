package service

import (
	"context"
	"sync"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// 单例模式
var postSrvIns *PostSrv
var postSrvOnce sync.Once

type PostSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetPostSrv() *PostSrv {
	postSrvOnce.Do(func() {
		postSrvIns = &PostSrv{}
	})
	return postSrvIns
}

// 重置单例 便于测试
func ResetPostSrv() {
	postSrvOnce = sync.Once{}
	postSrvIns = nil
}

// 创建帖子业务
func (s *PostSrv) CreatePost(ctx context.Context, req *types.PostCreateReq) (err error) {

	// 从c 取到当前发请求的用户userid
	// userID, err := request.GetLoginUserID(c)
	user, err := ctl.GetLoginUserID(ctx)
	if err != nil {
		zap.L().Error("ctl GetLoginUserID failed", zap.Error(err))
		// e.ResponseError(ctx, e.CodeNeedLogin)
		return err
	}
	// 获取发布帖子的用户id
	req.AuthorID = user.UserId
	postDao := mysql.NewPostDao(ctx)
	// 1. 生成post_id
	req.PostID = int64(snowflake.GenID())
	// 2. 保存到数据库
	// 参数绑定
	post := &models.Post{
		Status:      settings.POST_STATUS_ON,
		PostID:      req.PostID,
		CommunityID: req.CommunityID,
		AuthorID:    req.AuthorID,
		Title:       req.Title,
		Content:     req.Content,
	}
	err = postDao.CreatePost(post)
	if err != nil {
		zap.L().Error("mysql CreatePost failed", zap.Error(err))
		return
	}
	//3.  需要在redis记录帖子的创建时间
	err = redisCache.CreatePost(req.PostID, req.CommunityID)
	if err != nil {
		zap.L().Error("redisCache CreatePost failed", zap.Error(err))
		return
	}
	// 3. 返回
	return err
}

// 查询帖子列表业务
func (s *PostSrv) GetPostList(ctx context.Context, pageNum, pageSize int64) (postList []*types.PostListResp, err error) {
	postDao := mysql.NewPostDao(ctx)
	userDao := mysql.NewUserDao(ctx)
	communityDao := mysql.NewCommunityDao(ctx)
	// 查询数据库，找到所有的post 并返回分页的内容
	posts, err := postDao.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("mysql GetPostList failed", zap.Error(err))
		return nil, err
	}
	resp := make([]*types.PostListResp, 0, len(posts)) // 初始化内存空间
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := userDao.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := communityDao.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql GetCommunityDetailList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &types.PostListResp{
			AuthorName: user.UserName,
			// VoteAgreeNum: ,
			Post:            post,
			CommunityDetail: community,
		}
		// fmt.Println(postDetail.Post)
		// fmt.Println(postDetail.CommunityDetail)
		resp = append(resp, postDetail)
	}
	return resp, nil
}

// 查询帖子详情业务
func (s *PostSrv) GetPostDetailList(ctx context.Context, req *types.PostIdReq) (data *types.PostListResp, err error) {
	// 查询数据库，找到post_id的信息 并返回
	// 查询并组合接口需要的数据
	postDao := mysql.NewPostDao(ctx)
	post, err := postDao.GetPostDetailList(req.PostID)
	if err != nil {
		zap.L().Error("mysql GetPostDetailList failed",
			zap.Int64("pid", req.PostID),
			zap.Error(err))
		return
	}
	// fmt.Println(post)
	// fmt.Println(post.AuthorID)
	// 根据作者id 查询作者信息
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id 查询社区详情信息
	communityDao := mysql.NewCommunityDao(ctx)
	communityDetail, err := communityDao.GetCommunityDetailList(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql GetCommunityDetailList failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	data = &types.PostListResp{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return data, err
}

// 新版查询帖子业务，根据前端返回参数，创建时间或者分数进行排序
func (s *PostSrv) SearchPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. redis查询帖子id列表
	ids, err := redisCache.SearchPostIDsByOrder(p)
	if err != nil {
		zap.L().Error("redisCache SearchPostIDsByOrder failed",
			zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redisCache SearchPostIDsByOrder success, but empty")
		return
	}
	// 3. 根据id去数据库查询帖子详细信息
	// **返回的数据是按照给定的顺序
	posts, err := mysql.SearchPostListByIDs(ids)
	data = make([]*models.ApiPostDetail, 0, len(posts)) // 初始化内存空间
	if err != nil {
		zap.L().Error("mysql SearchPostListByIDs failed",
			zap.Error(err))
		return
	}

	// 提前查询好每条帖子的赞同数
	voteAgreeData, err := redisCache.GetPostAgreeByIDs(ids)
	if err != nil {
		return
	}

	// 将贴子的作者和社区信息查询处理填充到帖子中
	for index, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql GetCommunityDetailList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteAgreeNum:    voteAgreeData[index],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

// 根据社区id查询帖子的列表信息
func (s *PostSrv) CommunitySearchPostList(c context.Context, req *types.PostListReq) (data []*types.PostListResp, err error) {
	userDao := mysql.NewUserDao(c)
	// 2. redis查询帖子id列表
	ids, err := redisCache.CommunityPostIDsByOrder(req)
	if err != nil {
		zap.L().Error("redisCache CommunityPostIDsByOrder failed",
			zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redisCache CommunityPostIDsByOrder success, but empty")
		return
	}
	// 3. 根据id去数据库查询帖子详细信息
	// **返回的数据是按照给定的顺序

	posts, err := mysql.SearchPostListByIDs(ids)
	data = make([]*models.ApiPostDetail, 0, len(posts)) // 初始化内存空间
	if err != nil {
		zap.L().Error("mysql SearchPostListByIDs failed",
			zap.Error(err))
		return
	}

	// 提前查询好每条帖子的赞同数
	voteAgreeData, err := redisCache.GetPostAgreeByIDs(ids)
	if err != nil {
		return
	}

	// 将贴子的作者和社区信息查询处理填充到帖子中
	for index, post := range posts {
		// 根据作者id查询作者信息
		user, err := userDao.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql GetUserByID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailList(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql GetCommunityDetailList failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteAgreeNum:    voteAgreeData[index],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

// 封装整合查询帖子的业务 将SearchPostList和CommunitySearchPostList结合起来
func (s *PostSrv) GetPostListNew(c context.Context, req *types.PostListReq) (data []*types.PostListResp, err error) {
	if req.CommunityID == 0 {
		// 查询所有帖子
		data, err = s.SearchPostList(req)

	} else {
		// 按照社区id查询帖子
		data, err = s.CommunitySearchPostList(req)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
