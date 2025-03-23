package service

import (
	"context"
	"sync"

	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// 单例模式
var communitySrvIns *CommunitySrv
var communitySrvOnce sync.Once

type CommunitySrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetCommunitySrv() *CommunitySrv {
	communitySrvOnce.Do(func() {
		communitySrvIns = &CommunitySrv{}
	})
	return communitySrvIns
}

// 重置单例 便于测试
func ResetCommunitySrv() {
	communitySrvOnce = sync.Once{}
	communitySrvIns = nil
}

// 社区请求相关
func (s *CommunitySrv) GetCommunityList(ctx context.Context, req *types.CommuntityListReq) (resp interface{}, err error) {
	// 查询数据库，找到所有的community 并返回
	communityDao := mysql.NewCommunityDao(ctx)
	communityList, err := communityDao.GetCommunityList()
	if err != nil {
		zap.L().Error("mysql GetCommunityList failed", zap.Error(err))
		return nil, err
	}
	return communityList, nil
	// return mysql.GetCommunityList()
}

// 社区分类详情 根据community_id查询
func (s *CommunitySrv) GetCommunityDetailList(ctx context.Context, req *types.CommunityIdReq) (*models.Community, error) {
	communityDao := mysql.NewCommunityDao(ctx)
	communityList, err := communityDao.GetCommunityDetailList(req.CommunityID)
	if err != nil {
		zap.L().Error("mysql GetCommunityDetailList failed", zap.Error(err))
		return nil, err
	}
	// 查询数据库，找到community_id的信息 并返回
	// return mysql.GetCommunityDetailList(cid)
	return communityList, nil
}
