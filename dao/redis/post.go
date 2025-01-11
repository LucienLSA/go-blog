package redisCache

import (
	"github.com/LucienLSA/go-blog/models"
	"github.com/redis/go-redis/v9"
)

// 根据参数查询redis中id
func SearchPostIDsByOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求携带的order参数确定查询redis的key
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	// 第一页从0开始到pagesize，后续查询到start+pagesize-1
	start := (p.PageNum - 1) * p.PageSize
	end := start + p.PageSize - 1
	// 3. ZRevRange 按分数从大到小的查询ZSET中指定数量的key
	return rdb.ZRevRange(rctx, key, start, end).Result()
}

// 根据ids查询每条帖子的投赞成票数据
func GetPostAgreeByIDs(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := GetRedisKey(KeyPostVotedZSetPrefix + id)
	// 	// 查询1这个数据的key 统计帖子的赞成票
	// 	v := rdb.ZCount(rctx, key, "1", "1").Val()
	// 	data = append(data, v)
	// }
	// 使用pipeline一次发送多条命令，减少RTT
	zcountPip := rdb.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedZSetPrefix + id)
		zcountPip.ZCount(rctx, key, "1", "1")
	}
	cmders, err := zcountPip.Exec(rctx)
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return data, err
}
