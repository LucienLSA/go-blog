package types

// 投票数据
type PostVoteDataReq struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id,string" binding:"required"`  // 投票帖子id
	Kind   int8  `json:"kind,string" binding:"oneof=1 0 -1"` // 投票帖子赞成(1)\反对(-1)\取消投票(0)
}
