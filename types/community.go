package types

type CommuntityListReq struct {
}

// 社区分类查询 id
type CommunityIdReq struct {
	CommunityID int64 `json:"community_id" form:"community_id" example:"2" binding:"required"` // 社区id
}
