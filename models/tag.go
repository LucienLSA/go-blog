package model

type Tag struct {
	*Model
	ID    int64  `json:"id" gorm:"column:tag_id; primary_key;not null;auto_increment" comment:"标签ID"`
	Name  string `json:"name" gorm:"column:tag_name; varchar(100) ;default:''" comment:"标签名称"`
	State uint8  `json:"state" gorm:"column:tag_state; tinyint(3) ;default:1 comment:'状态 0-禁用 1-启用"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}
