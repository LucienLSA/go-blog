package service

import (
	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/models"
)

// 社区请求相关
func GetCommunityList() ([]*models.Community, error) {
	// 查询数据库，找到所有的community 并返回
	return mysql.GetCommunityList()
}
