package service

import (
	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
)

// 存放业务逻辑

// 注册业务
func SignUp(p *models.ParamSignUp) {
	// 1. 判断用户是否存在
	mysql.QueryUserByUsername()
	// 2. 生成UID
	snowflake.GenID()

	// 3. 保存到数据库
	mysql.InsertUser()
}
