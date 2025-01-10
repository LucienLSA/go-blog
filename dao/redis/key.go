package redisCache

// key健 定义常量易于查询和拆分
// 冒号分割命名空间 使用其区分不同的key
const (
	Prefix                 = "bluebell:"
	KeyPostTimeZSet        = "post:time"   // zset; 帖子的发帖时间
	KeyPostScoreZSet       = "post:score"  // zset; 帖子的投票分数
	KeyPostVotedZSetPrefix = "post:voted:" // zset; 记录用户及投票类型；参数是post_id
)

// 给redis key加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}
