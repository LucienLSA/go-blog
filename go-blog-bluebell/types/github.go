package types

// GitHubTrendingData GitHub热点数据结构
type GitHubTrendingData struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GitHubTrendingRequest GitHub热点数据请求结构
type GitHubTrendingRequest struct {
	Language string `json:"language" form:"language"` // 编程语言筛选
	Since    string `json:"since" form:"since"`       // 时间范围：daily, weekly, monthly
}

// GitHubTrendingResponse GitHub热点数据响应结构
type GitHubTrendingResponse struct {
	Language string                `json:"language"`
	Since    string                `json:"since"`
	Data     []GitHubTrendingData `json:"data"`
	Total    int                   `json:"total"`
} 