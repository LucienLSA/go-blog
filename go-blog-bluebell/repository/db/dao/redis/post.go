package redisCache

import (
	"strconv"
	"time"

	"github.com/LucienLSA/go-blog/types"
	"github.com/redis/go-redis/v9"
)

// 根据参数查询redis中id
func SearchPostIDsByOrder(req *types.PostListReq) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求携带的order参数确定查询redis的key
	key := GetRedisKey(KeyPostTimeZSet)
	if req.Order == types.OrderScore {
		key = GetRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	// 第一页从0开始到pagesize，后续查询到start+pagesize-1
	// 3. ZRevRange 按分数从大到小的查询ZSET中指定数量的key
	return getIDsFromKey(key, req.PageNum, req.PageSize)
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

// 根据社区id参数查询redis中ids
func CommunityPostIDsByOrder(req *types.PostListReq) ([]string, error) {
	// 1. 确定排序依据的ZSet（时间或分数）
	orderKey := GetRedisKey(KeyPostTimeZSet)
	if req.Order == types.OrderScore {
		orderKey = GetRedisKey(KeyPostScoreZSet)
	}

	// 2. 定义社区帖子集合的Key和临时合并结果的Key
	cKey := GetRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(req.CommunityID))) // 社区下的帖子ID集合（Set）
	key := orderKey + strconv.Itoa(int(req.CommunityID))                            // 临时合并结果的ZSet Key

	// 3. 检查临时ZSet是否存在，不存在则通过ZInterStore生成
	if rdb.Exists(rctx, key).Val() < 1 {
		// 使用管道批量执行命令
		pipeline := rdb.Pipeline()
		// ZInterStore合并社区帖子集合和排序ZSet，取MAX分数作为合并后分数
		pipeline.ZInterStore(rctx, key, &redis.ZStore{
			Keys:      []string{cKey, orderKey}, // 参与合并的两个集合：社区帖子Set + 排序ZSet
			Aggregate: "MAX",                    // 合并策略：取两个集合中相同元素的最大分数（保留排序权重）
		})
		// 设置临时ZSet的过期时间为60秒，减少重复计算
		pipeline.Expire(rctx, key, 60*time.Second)
		// 执行管道命令
		_, err := pipeline.Exec(rctx)
		if err != nil {
			return nil, err
		}
	}

	// 4. 从临时ZSet中分页查询帖子ID
	return getIDsFromKey(key, req.PageNum, req.PageSize)
}
