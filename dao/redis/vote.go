package redisCache

import (
	"errors"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

// 基于用户投票的相关算法

// 简化版的投票分数算法
// 投一票加432分， 随着时间越新，新的帖子分数越高； 一天有86400秒，86400/432=200表示一天有200张赞同票能使得帖子在首页

// 投票的几种情况
// kind为1有两种情况
// 1. 之前没有投过票，现在投赞成票 --> 更新分数和投票记录 差值绝对值1 +432
// 2. 之前投反对票，现在改投赞成票 --> 更新分数和投票记录 差值绝对值2 +432*2
// kind为0有两种情况
// 1. 之前没有投过赞成票，现在取消投票 --> 更新分数和投票记录 差值绝对值1 -432
// 2. 之前没有投过反对票，现在取消投票 --> 更新分数和投票记录 差值绝对值1 +432
// kind为-1有两种情况
// 1. 之前没有投过票，现在投反对票 --> 更新分数和投票记录 差值绝对值1 -432
// 2. 之前投赞成票，现在改投反对票 --> 更新分数和投票记录 差值绝对值2 -432*2
// 投票限制
// 每个帖子自发表日起一个星期之内允许用户投票，超过一个星期就不允许投票了
// 1. 到期之后，将redis中保存的赞成票和反对票存储到mysql中，进行持久化存储
// 2. 到期之后，删除KeyPostVotedZSetPrefix

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePreVote     = 432 // 每一票多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID int64) (err error) {
	// 事务操作
	pipeline := rdb.TxPipeline()
	// 帖子时间
	rdb.ZAdd(rctx, GetRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数
	rdb.ZAdd(rctx, GetRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err = pipeline.Exec(rctx)
	// fmt.Println(err)
	return err
}

// 为帖子投票功能
func PostVote(uid, pid string, kind float64) (err error) {
	// rdb := redisCache.RedisClient
	// 1. 判断投票限制
	// 取redis帖子的发布时间
	postTime := rdb.ZScore(rctx, GetRedisKey(KeyPostTimeZSet), pid).Val() // 返回float64
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2，更新帖子分数
	// 先查当前用户给当前帖子之前的投票记录
	// 定义给帖子投票的记录
	keyPostVoted := KeyPostVotedZSetPrefix + pid
	userPreVote := rdb.ZScore(rctx, GetRedisKey(keyPostVoted), uid).Val()

	// 如果本次投票和之前的投票类型是一致的，则不允许重复投票
	if kind == userPreVote {
		return ErrVoteRepeated
	}

	var k float64
	if kind > userPreVote {
		k = 1
	} else {
		k = -1
	}
	// 计算分数值
	diff := math.Abs(userPreVote - kind) // 计算两次投票的差值，之前投票-当前投票
	// 对帖子增加分数
	txpipline := rdb.TxPipeline()
	calScore := k * diff * scorePreVote
	txpipline.ZIncrBy(rctx, GetRedisKey(KeyPostScoreZSet), calScore, pid)
	// if ErrVoteTimeExpire == nil {
	// 	return err
	// }

	// 3. 记录用户为该帖子投过票
	if kind == 0 {
		// 判断是否取消投票，是则删除该记录
		txpipline.ZRem(rctx, GetRedisKey(keyPostVoted), pid)
	} else {
		// 记录投票
		// redis.Z为特定的结构体
		txpipline.ZAdd(rctx, GetRedisKey(keyPostVoted), redis.Z{
			Score:  kind, // 赞成票或反对
			Member: uid,  // 哪个用户
		})
	}
	// 执行命令
	_, err = txpipline.Exec(rctx)
	return err
}
