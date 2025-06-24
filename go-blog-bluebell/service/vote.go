package service

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// 单例模式
var voteSrvIns *VoteSrv
var voteSrvOnce sync.Once

type VoteSrv struct {
}

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
func GetVoteSrv() *VoteSrv {
	voteSrvOnce.Do(func() {
		voteSrvIns = &VoteSrv{}
	})
	return voteSrvIns
}

// 重置单例 便于测试
func ResetVoteSrv() {
	userSrvOnce = sync.Once{}
	userSrvIns = nil
}

// PostVote 为帖子投票
func (s *VoteSrv) PostVote(ctx context.Context, req *types.PostVoteDataReq) (err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil || u == nil {
		return errors.New("用户未登录")
	}
	uid := u.UserId
	zap.L().Debug("PostVote", zap.Int64("userID", uid),
		zap.Int64("PostID", req.PostID),
		zap.Int8("Kind", req.Kind))
	err = redisCache.PostVote(strconv.Itoa(int(uid)),
		strconv.FormatInt(req.PostID, 10), float64(req.Kind))
	return err
}
