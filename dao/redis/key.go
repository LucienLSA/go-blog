package redisCache

// key健 定义常量易于查询和拆分
// 冒号分割命名空间 使用其区分不同的key
const (
	Prefix                 = "bluebell:"
	KeyPostTimeZSet        = "post:time"   // zset; 帖子的发帖时间
	KeyPostScoreZSet       = "post:score"  // zset; 帖子的投票分数
	KeyPostVotedZSetPrefix = "post:voted:" // zset; 记录用户及投票类型；参数是post_id

	KeyCommunitySetPrefix = "community:" // set; 保存社区下帖子的id
)

// 给redis key加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}

// 2. 确定查询的索引起始点
// 第一页从0开始到pagesize，后续查询到start+pagesize-1
// 3. ZRevRange 按分数从大到小的查询ZSET中指定数量的key
func getIDsFromKey(key string, pageNum, pageSize int64) ([]string, error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize - 1
	return rdb.ZRevRange(rctx, key, start, end).Result()
}
