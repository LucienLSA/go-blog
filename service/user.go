package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/dao/redis"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/email"
	"github.com/LucienLSA/go-blog/pkg/jwt"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/pkg/upload"
	"github.com/LucienLSA/go-blog/settings"
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
	// 3. 构造一个user示例
	// 获取默认头像
	defaultAvatar := settings.Conf.AppConfig.PhotoPathConfig
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Age:      p.Age,
		Email:    p.Email,
		Password: p.Password,
		Gender:   p.Gender,
		Avatar:   defaultAvatar.DefaultAvatar,
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

// 用户信息修改
func Update(uid int64, p *models.ParamUpdate) (err error) {
	var user *models.User
	// 1. 查询登录的用户ID，不用在数据库中查询是否存在，因为是已经登录过的
	user, err = mysql.GetUserByID(uid)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}
	user = &models.User{
		UserID:   uid,
		Username: p.Username,
		Age:      p.Age,
		Email:    p.Email,
		Password: p.Password,
		Gender:   p.Gender,
	}
	err = mysql.UpdateUser(uid, user)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return err
	}
	return err
}

// 上传用户头像
func UploadAvatar(uId int64, file multipart.File, fileSize int64) (paramAvatar *models.ParamAvatar, err error) {
	var user *models.User
	user, err = mysql.GetUserByID(uId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return nil, err
	}
	// 保存到本地
	var path string
	if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal { // 兼容两种存储方式
		path, err = upload.UploadAvatarToLocalStatic(file, uId, user.Username)
	} else { //保存到七牛云oss
		path, err = upload.UploadToQiNiuAvatar(file, user.Username, fileSize)
	}

	if err != nil {
		zap.L().Error("upload avatar failed", zap.Error(err))
		return nil, err
	}
	user.Avatar = path
	err = mysql.UpdateAvatar(uId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return nil, err
	}
	// fmt.Println(user)
	paramAvatar = &models.ParamAvatar{
		UserID:   uId,
		UserName: user.Username,
		Avatar:   upload.AvatarURL() + user.Avatar,
	}
	// if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal {
	// 	paramAvatar.Avatar = upload.AvatarURL() + user.Avatar
	// } else {
	// 	paramAvatar = &models.ParamAvatar{
	// 		Avatar:   upload.AvatarURL() + user.Avatar,
	// 		UserID:   uId,
	// 		UserName: user.Username,
	// 	}
	// }

	return paramAvatar, nil
}

// 用户发送邮箱验证码
func SendEmail(uid int64, p *models.ParamSendEmail) (err error) {
	var user *models.User
	user, err = mysql.GetUserByID(uid)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}
	var address string
	// notice := new(models.Notice)
	var notice *models.Notice
	token, err := jwt.GenerateEmailToken(p.OperationType, uid, p.Email, user.Password)
	if err != nil {
		zap.L().Error("jwt GenerateEmailToken failed", zap.Error(err))
		return err
	}
	notice, err = mysql.GetNoticeById(p.OperationType)
	if err != nil {
		zap.L().Error("mysql GetNoticeById failed", zap.Error(err))
		return err
	}
	address = settings.Conf.EmailConfig.VaildEmail + token
	mailStr := notice.Text
	mailTex := strings.Replace(mailStr, "Email", address, -1)
	err = email.Send(mailTex, p.Email, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
	if err != nil {
		zap.L().Error("email Send failed", zap.Error(err))
		return err
	}
	// m := mail.NewMessage()
	// m.SetHeader("From", settings.Conf.EmailConfig.SmtpEmail)
	// m.SetHeader("To", p.Email)
	// m.SetHeader("Subject", "bluebell")
	// m.SetBody("text/html", mailTex)
	// d := mail.NewDialer(settings.Conf.EmailConfig.SmtpHost, 465, settings.Conf.EmailConfig.SmtpEmail, settings.Conf.EmailConfig.SmtpPass)
	// d.StartTLSPolicy = mail.MandatoryStartTLS
	// if err = d.DialAndSend(m); err != nil {
	// 	zap.L().Error("DialAndSend failed", zap.Error(err))
	// 	return err
	// }
	return
}

func ValidEmail(c context.Context, token string) (err error) {
	var operationType int
	var uId int64
	var user *models.User
	var email string
	// 1.校验参数
	claims, err := jwt.ParseEmailToken(token)
	if err != nil {
		zap.L().Error("jwt ParseEmailToken failed", zap.Error(err))
		return err
	}
	if time.Now().Unix() > claims.ExpiresAt {
		zap.L().Error("check token timeout")
		return err
	}
	uId = claims.UserID
	email = claims.Email
	operationType = claims.OperationType
	user, err = mysql.GetUserByID(uId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}
	if operationType == 1 {
		user.Email = email
	} else if operationType == 2 {
		user.Email = " "
	}
	err = mysql.UpdateUserEmail(uId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUserEmail failed", zap.Error(err))
		return err
	}
	return
}
