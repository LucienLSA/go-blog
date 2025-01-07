package service

import (
	"fmt"

	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/jwt"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
)

// 存放业务逻辑

// 注册业务
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询错误
		return err
	}

	// 2. 生成UID
	userID := snowflake.GenID()
	fmt.Println(userID)
	// 构造一个user示例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Age:      p.Age,
		Email:    p.Email,
		Password: p.Password,
		Gender:   p.Gender,
	}
	fmt.Println(&user)
	// 3. 保存到数据库
	err = mysql.InsertUser(user)
	return err
}

// 登录业务
func Login(p *models.ParamLogin) (token string, err error) {
	//  登录
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 将用户登录输入的名称和密码信息传入dao层
	if err = mysql.Login(user); err != nil {
		// 传递的是指针，能拿到数据库中注册时原本生成的UserID和Username
		return "", err
	}
	// 生成JWT
	token, err = jwt.GenToken(user.UserID, user.Username)
	return token, err
}
