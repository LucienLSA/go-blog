package model

type Article struct {
	*Model
	ID              int64  `json:"id" gorm:"column:article_id;primary_key;auto_increment;not null"`                                 // 文章ID
	Title           string `json:"title" gorm:"column:article_title;varchar(100) ;default:''" comment:"文章标题"`                       // 文章标题
	Desc            string `json:"desc" gorm:"column:article_desc;varchar(255) ;default:''" comment:"文章简介"`                         // 文章简介
	Cover_image_url string `json:"cover_image_url" gorm:"column:article_cover_image_url;varchar(255) ;default:''" comment:"封面图片地址"` // 文章封面图片地址
	Content         string `json:"content" gorm:"column:article_content" comment:"文章内容"`
	State           uint8  `json:"state" gorm:"column:article_state;tinyint(3) ;default:1" comment:"状态 0-禁用 1-启用"` // 状态 0-禁用 1-启用
	// 创建文章标签关联表，这个表主要用于记录文章和标签之间的 1:N 的关联关系
	Article_Tag []Tag
}

func (a Article) TableName() string {
	return "blog_article"
}
