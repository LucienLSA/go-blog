package types

type CommuntityListReq struct {
}

// 社区分类查询 id
type CommunityIdReq struct {
	// CommunityID int64 `json:"community_id" form:"community_id" example:"2"` // 社区id
}

type CommunityCreateResp struct {
	CommunityID   int64  ` json:"community_id" form:"community_id"`     // 社区id
	CommunityName string ` json:"community_name" form:"community_name"` // 社区名称
	Introduction  string ` json:"introduction" form:"introduction"`     // 社区介绍
}

type CommunityDetailResp struct {
	CommunityID   int64  ` json:"community_id" form:"community_id"`     // 社区id
	CommunityName string ` json:"community_name" form:"community_name"` // 社区名称
	Introduction  string ` json:"introduction" form:"introduction"`     // 社区介绍
}
