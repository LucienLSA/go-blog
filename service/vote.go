package service

import (
	"strconv"

	redisCache "github.com/LucienLSA/go-blog/dao/redis"
	"github.com/LucienLSA/go-blog/models"
	"go.uber.org/zap"
)

// PostVote 为帖子投票
func PostVote(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("PostVote", zap.Int64("userID", userID),
		zap.Int64("PostID", p.PostID),
		zap.Int8("Kind", p.Kind))
	return redisCache.PostVote(strconv.FormatInt(userID, 10),
		strconv.FormatInt(p.PostID, 10), float64(p.Kind))
}
