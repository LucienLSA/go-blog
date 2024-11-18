package model

// 手动建立关联表
type Author struct {
	*Model
	ID       int64 `json:"id" gorm:"column:id;primary_key;not null;auto_increment"`
	UserName int64 `json:"user_name" gorm:"column:user_name;default:'';varchar(50)" comment:"账号"`
	Password int64 `json:"password" gorm:"column:password;default:'';varchar(50)" comment:"密码"`
}

func (Author) TableName() string {
	return "blog_auth"
}
