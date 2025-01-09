package service

import (
	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"go.uber.org/zap"
)

// 创建帖子业务
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post_id
	p.PostID = int64(snowflake.GenID())
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	// 3. 返回
	return
}

// 查询帖子列表业务
func GetPost() ([]*models.PostList, error) {
	// 查询数据库，找到所有的post 并返回
	return mysql.GetPost()
}

// 查询帖子详情业务
func GetPostDetailList(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询数据库，找到post_id的信息 并返回
	// 查询并组合接口需要的数据
	post, err := mysql.GetPostDetailList(pid)
	if err != nil {
		zap.L().Error("mysql GetPostDetailList failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// fmt.Println(post)
	// fmt.Println(post.AuthorID)
	// 根据作者id 查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据社区id 查询社区详情信息
	communityDetail, err := mysql.GetCommunityDetailList(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql GetCommunityDetailList failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return data, err
}
