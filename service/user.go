package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/LucienLSA/go-blog/pkg/email"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/settings"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// 单例模式
var userSrvIns *UserSrv
var userSrvOnce sync.Once

type UserSrv struct {
}

// 单例实例 不对外暴露
func GetUserSrv() *UserSrv {
	userSrvOnce.Do(func() {
		userSrvIns = &UserSrv{}
	})
	return userSrvIns
}

// 重置单例 便于测试
func ResetUserSrv() {
	userSrvOnce = sync.Once{}
	userSrvIns = nil
}

// 存放业务逻辑

// 注册业务
func (s *UserSrv) UserSignUp(ctx context.Context, req *types.UserSignUpReq) (err error) {
	userDao := dao.NewUserDao(ctx)
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询错误
		zap.L().Error("mysql CheckUserExist failed", zap.Error(err))
		return err
	}
	// 判断邮箱格式是否正确
	if err = email.VerifyEmailFormat(p.Email); err != nil {
		// 邮箱格式不正确
		zap.L().Error("email VerifyEmailFormat failed", zap.Error(err))
		return err
	}

	// 判断邮箱是否重复注册
	if err = mysql.ExistUserEmail(p.Email); err != nil {
		// 数据库查询错误
		zap.L().Error("mysql ExistUserEmai failed", zap.Error(err))
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

// // 登录业务
// func Login(p *models.ParamLogin) (user *models.User, err error) {
// 	//  登录
// 	user = &models.User{
// 		Username: p.Username,
// 		Password: p.Password,
// 	}
// 	// 将用户登录输入的名称和密码信息传入dao层
// 	if err = mysql.Login(user); err != nil {
// 		// 传递的是指针，能拿到数据库中注册时原本生成的UserID和Username
// 		zap.L().Error("mysql Login failed", zap.Error(err))
// 		return nil, err
// 	}
// 	// 生成JWT
// 	token, err := jwt.GenToken(user.UserID, user.Username)
// 	if err != nil {
// 		zap.L().Error("jwt GenToken failed", zap.Error(err))
// 	}
// 	user.Token = token
// 	// 保存到redis中
// 	if err = redisCache.StorgeUserIdToken(token, user.Username); err != nil {
// 		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
// 	}
// 	return
// }

// // 登录发送邮箱验证码业务
// func SendEmailCode(p *models.ParamSendEmailCode) (err error) {
// 	// 检查是否在1分钟内发送过邮件， 如果有发送过，需等待后才能发送
// 	err = redisCache.EmailCodeExists(p.UserEmail)
// 	if err != nil {
// 		zap.L().Error("send email code from redis in 1 minutes")
// 		fmt.Printf("send email code from redis in 1 minutes")
// 		return err
// 	}

// 	// 获取六位数邮箱验证码
// 	code := email.GetConfirmCode()

// 	// 将其存储至Redis中，由于Redis为KV键值对存储所以需要定义前缀方便使用
// 	err = redisCache.StorgeEmailCode(p.UserEmail, code)
// 	if err != nil {
// 		zap.L().Error("redisCache.StorgeEmailCode failed, err:%\v", zap.Error(err))
// 		fmt.Printf("redisCache.StorgeEmailCode failed, err:%\v", err)
// 		return err
// 	}

// 	// 发送邮件，此处为方便起见没有处理返回值
// 	// address = settings.Conf.EmailConfig.VaildEmail + token
// 	var notice *models.Notice
// 	notice, err = mysql.GetNoticeById(p.OperationType)
// 	if err != nil {
// 		zap.L().Error("mysql GetNoticeById failed", zap.Error(err))
// 		return err
// 	}
// 	mailStr := notice.Text
// 	var builder strings.Builder
// 	builder.WriteString(mailStr)
// 	builder.Write([]byte(code))
// 	// mailTex := strings.Replace(mailStr, "Email", code, -1)
// 	// 第一个参数内容 第二个参数接收方 第三个标题 第四个发送方
// 	err = email.Send(builder.String(), p.UserEmail, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
// 	if err != nil {
// 		zap.L().Error("email Send failed", zap.Error(err))
// 		return err
// 	}
// 	// 设置每个邮箱发送邮件的时间 此处设置为1分钟，由于Redis为KV键值对存储所以需要定义前缀方便使用
// 	err = redisCache.SendEmailCode(p.UserEmail, code)
// 	if err != nil {
// 		zap.L().Error("redisCache.SendEmailCode failed, err:%\v", zap.Error(err))
// 		fmt.Printf("redisCache.SendEmailCode failed, err:%\v", err)
// 		return err
// 	}
// 	return nil
// }

// // 用户邮箱验证码登录业务
// func LoginEmail(p *models.ParamLoginEmail) (user *models.User, err error) {
// 	// 将用户邮箱登录输入的邮箱传入dao层mysql中进行判断是否存在，并取出
// 	// 从数据库中获取用户信息（用户名和id）
// 	user, err = mysql.GetUserByEmail(p.UserEmail)
// 	// 如果不存在则错误
// 	if err != nil {
// 		zap.L().Error(" mysql.NotExistUserEmail failed, err:%\v", zap.Error(err))
// 		fmt.Printf(" mysql.NotExistUserEmail failed, err:%\v", err)
// 		return nil, err
// 	}
// 	userInfo := &models.User{
// 		UserID:   user.ID,
// 		Email:    p.UserEmail,
// 		Username: user.Username,
// 	}
// 	// 校验验证码正确性
// 	err = redisCache.CheckEmailCode(p.UserEmail, userInfo.Email, p.Code)
// 	if err != nil {
// 		zap.L().Error("redisCache.CheckEmailCode failed, err:%\v", zap.Error(err))
// 		fmt.Printf("redisCache.CheckEmailCode failed, err:%\v", err)
// 		return nil, err
// 	}
// 	// 生成JWT
// 	token, err := jwt.GenToken(userInfo.UserID, userInfo.Username)
// 	if err != nil {
// 		zap.L().Error("jwt GenToken failed", zap.Error(err))
// 	}
// 	user.Token = token
// 	// 保存到redis中
// 	if err = redisCache.StorgeUserIdToken(token, userInfo.Email); err != nil {
// 		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
// 	}
// 	return
// }

// // 用户信息修改
// func Update(uid uint, p *models.ParamUpdate) (err error) {
// 	var user *models.User
// 	// 1. 查询登录的用户ID，不用在数据库中查询是否存在，因为是已经登录过的
// 	user, err = mysql.GetUserByID(uid)
// 	if err != nil {
// 		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
// 		return err
// 	}
// 	user = &models.User{
// 		UserID:   uid,
// 		Username: p.Username,
// 		Age:      p.Age,
// 		Email:    p.Email,
// 		Password: p.Password,
// 		Gender:   p.Gender,
// 	}
// 	err = mysql.UpdateUser(uid, user)
// 	if err != nil {
// 		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
// 		return err
// 	}
// 	return err
// }

// // 上传用户头像
// func UploadAvatar(uId int64, file multipart.File, fileSize int64) (paramAvatar *models.ParamAvatar, err error) {
// 	var user *models.User
// 	user, err = mysql.GetUserByID(uId)
// 	if err != nil {
// 		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
// 		return nil, err
// 	}
// 	// 保存到本地
// 	var path string
// 	if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal { // 兼容两种存储方式
// 		path, err = upload.UploadAvatarToLocalStatic(file, uId, user.Username)
// 	} else { //保存到七牛云oss
// 		path, err = upload.UploadToQiNiuAvatar(file, user.Username, fileSize)
// 	}

// 	if err != nil {
// 		zap.L().Error("upload avatar failed", zap.Error(err))
// 		return nil, err
// 	}
// 	user.Avatar = path
// 	err = mysql.UpdateAvatar(uId, user)
// 	if err != nil {
// 		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
// 		return nil, err
// 	}
// 	// fmt.Println(user)
// 	paramAvatar = &models.ParamAvatar{
// 		UserID:   uId,
// 		UserName: user.Username,
// 		Avatar:   upload.AvatarURL() + user.Avatar,
// 	}
// 	// if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal {
// 	// 	paramAvatar.Avatar = upload.AvatarURL() + user.Avatar
// 	// } else {
// 	// 	paramAvatar = &models.ParamAvatar{
// 	// 		Avatar:   upload.AvatarURL() + user.Avatar,
// 	// 		UserID:   uId,
// 	// 		UserName: user.Username,
// 	// 	}
// 	// }

// 	return paramAvatar, nil
// }

// // 用户发送邮箱验证码
// func SendEmail(uid int64, p *models.ParamSendEmail) (err error) {
// 	var user *models.User
// 	user, err = mysql.GetUserByID(uid)
// 	if err != nil {
// 		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
// 		return err
// 	}
// 	var address string
// 	// notice := new(models.Notice)
// 	var notice *models.Notice
// 	token, err := jwt.GenerateEmailToken(p.OperationType, uid, p.Email, user.Password)
// 	if err != nil {
// 		zap.L().Error("jwt GenerateEmailToken failed", zap.Error(err))
// 		return err
// 	}
// 	notice, err = mysql.GetNoticeById(p.OperationType)
// 	if err != nil {
// 		zap.L().Error("mysql GetNoticeById failed", zap.Error(err))
// 		return err
// 	}
// 	address = settings.Conf.EmailConfig.VaildEmail + token
// 	mailStr := notice.Text
// 	mailTex := strings.Replace(mailStr, "Email", address, -1)
// 	err = email.Send(mailTex, p.Email, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
// 	if err != nil {
// 		zap.L().Error("email Send failed", zap.Error(err))
// 		return err
// 	}
// 	// m := mail.NewMessage()
// 	// m.SetHeader("From", settings.Conf.EmailConfig.SmtpEmail)
// 	// m.SetHeader("To", p.Email)
// 	// m.SetHeader("Subject", "bluebell")
// 	// m.SetBody("text/html", mailTex)
// 	// d := mail.NewDialer(settings.Conf.EmailConfig.SmtpHost, 465, settings.Conf.EmailConfig.SmtpEmail, settings.Conf.EmailConfig.SmtpPass)
// 	// d.StartTLSPolicy = mail.MandatoryStartTLS
// 	// if err = d.DialAndSend(m); err != nil {
// 	// 	zap.L().Error("DialAndSend failed", zap.Error(err))
// 	// 	return err
// 	// }
// 	return
// }

// // 用户绑定邮箱和解绑邮箱业务
// func ValidEmail(c context.Context, token string) (err error) {
// 	var operationType int
// 	var uId int64
// 	var user *models.User
// 	var email string
// 	// 1.校验参数
// 	claims, err := jwt.ParseEmailToken(token)
// 	if err != nil {
// 		zap.L().Error("jwt ParseEmailToken failed", zap.Error(err))
// 		return err
// 	}
// 	if time.Now().Unix() > claims.ExpiresAt {
// 		zap.L().Error("check token timeout")
// 		return err
// 	}
// 	uId = claims.UserID
// 	email = claims.Email
// 	operationType = claims.OperationType
// 	user, err = mysql.GetUserByID(uId)
// 	if err != nil {
// 		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
// 		return err
// 	}
// 	if operationType == 1 {
// 		user.Email = email
// 	} else if operationType == 2 {
// 		user.Email = " "
// 	}
// 	err = mysql.UpdateUserEmail(uId, user)
// 	if err != nil {
// 		zap.L().Error("mysql UpdateUserEmail failed", zap.Error(err))
// 		return err
// 	}
// 	return
// }
