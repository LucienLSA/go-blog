package service

import (
	"fmt"
	"mime/multipart"

	"github.com/LucienLSA/go-blog/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/dao/redis"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/jwt"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/pkg/upload"
	"go.uber.org/zap"
)

// 存放业务逻辑

// 注册业务
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询错误
		zap.L().Error("mysql CheckUserExist failed", zap.Error(err))
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
	if err != nil {
		zap.L().Error("mysql InsertUser failed", zap.Error(err))
	}
	return err
}

// 登录业务
func Login(p *models.ParamLogin) (user *models.User, err error) {
	//  登录
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 将用户登录输入的名称和密码信息传入dao层
	if err = mysql.Login(user); err != nil {
		// 传递的是指针，能拿到数据库中注册时原本生成的UserID和Username
		zap.L().Error("mysql Login failed", zap.Error(err))
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		zap.L().Error("jwt GenToken failed", zap.Error(err))
	}
	user.Token = token
	// 保存到redis中
	if err = redisCache.StorgeUserIdToken(token, user.Username); err != nil {
		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
	}
	return
}

// 上传用户头像
func UploadAvatar(uId int64, file multipart.File, fileSize int64) (resp interface{}, err error) {
	var user *models.User
	user, err = mysql.GetUserByID(uId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return nil, err
	}
	// 保存到本地
	path, err := upload.UploadAvatarToLocalStatic(file, uId, user.Username)
	if err != nil {
		zap.L().Error("upload UploadAvatarToLocalStatic failed", zap.Error(err))
		return nil, err
	}
	user.Avatar = path
	err = mysql.UpdateUser(uId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return nil, err
	}
	resp = &models.ParamUserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Age:      user.Age,
	}
	return resp, nil
}
