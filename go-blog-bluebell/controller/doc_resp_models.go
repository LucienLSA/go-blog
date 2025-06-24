package controller

// 社区 swagger 展示结构体
// @Description 社区信息
// @name SwaggerCommunity
// @property community_id   int64  "社区ID"
// @property community_name string "社区名称"
// @property introduction   string "社区介绍"
type SwaggerCommunity struct {
	CommunityID   int64  `json:"community_id"`
	CommunityName string `json:"community_name"`
	Introduction  string `json:"introduction"`
}

// 帖子 swagger 展示结构体
// @Description 帖子信息
// @name SwaggerPost
// @property post_id      int64  "帖子ID"
// @property title        string "标题"
// @property content      string "内容"
// @property author_id    int64  "作者ID"
// @property community_id int64  "社区ID"
// @property status       int32  "状态"
type SwaggerPost struct {
	PostID      int64  `json:"post_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id"`
	Status      int32  `json:"status"`
}

// 用户 swagger 展示结构体
// @Description 用户信息
// @name SwaggerUser
// @property user_id   int64  "用户ID"
// @property user_name string "用户名"
// @property email     string "邮箱"
// @property avatar    string "头像"
type SwaggerUser struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// 通用成功响应
// @Description 通用成功响应
// @name _ResponseSuccess
// @property code    int         "状态码"
// @property message string      "消息"
// @property data    interface{} "数据"
type _ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 通用错误响应
// @Description 通用错误响应
// @name _ResponseError
// @property code    int    "状态码"
// @property message string "消息"
// @property data    string "数据"
type _ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// 社区列表响应
// @Description 社区列表响应
// @name _ResponseCommunityList
// @property code    int                "状态码"
// @property message string             "消息"
// @property data    []SwaggerCommunity "社区列表"
type _ResponseCommunityList struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []SwaggerCommunity `json:"data"`
}

// 社区详情响应
// @Description 社区详情响应
// @name _ResponseCommunityDetail
// @property code    int              "状态码"
// @property message string           "消息"
// @property data    SwaggerCommunity "社区详情"
type _ResponseCommunityDetail struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    SwaggerCommunity `json:"data"`
}

// 帖子列表响应
// @Description 帖子列表响应
// @name _ResponsePostList
// @property code    int           "状态码"
// @property message string        "消息"
// @property data    []SwaggerPost "帖子列表"
type _ResponsePostList struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []SwaggerPost `json:"data"`
}

// 帖子详情响应
// @Description 帖子详情响应
// @name _ResponsePostDetail
// @property code    int         "状态码"
// @property message string      "消息"
// @property data    SwaggerPost "帖子详情"
type _ResponsePostDetail struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    SwaggerPost `json:"data"`
}

// 用户登录响应
// @Description 用户登录响应
// @name _ResponseUserLogin
// @property code    int        "状态码"
// @property message string     "消息"
// @property data    SwaggerUser "用户信息"
type _ResponseUserLogin struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    SwaggerUser `json:"data"`
}
