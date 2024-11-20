package models

import (
	"blog-service/global"
)

// 手动建立关联表
// 作者信息
type Author struct {
	ID       int64  `json:"id" gorm:"column:id;primary_key;not null;auto_increment"`
	Username string `json:"username" gorm:"column:username;varchar(50);unique"`
	Password string `json:"password" gorm:"column:password;varchar(50)"`
}

func (Author) TableName() string {
	return "blog_auth"
}

// // 用户密码加密
// func SetPassword(password string) error {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	user
// }

// 检查用户名和密码是否匹配
func CheckAuthor(username, password string) bool {
	var author Author
	global.DBEngine.Select("id").
		Where(Author{Username: username, Password: password}).
		First(&author)
	if author.ID > 0 {
		return true
	}
	return false
}
