package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	"github.com/LucienLSA/go-blog/pkg/ctl"
	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/pkg/email"
	"github.com/LucienLSA/go-blog/pkg/jwt"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/pkg/upload"
	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
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

// 单例实例 不对外暴露 通过GetUserSrv来返回实例对象
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
	userDao := mysql.NewUserDao(ctx)
	// 1. 判断用户是否存在
	_, exist, err := userDao.CheckUserExist(req.UserName)
	if err != nil {
		// 数据库查询错误
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return err
	}
	if exist {
		// 用户存在
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return e.ErrorUserExist
	}
	// fmt.Println(req.Email)
	if req.Email != "" {
		// 判断邮箱格式是否正确
		if err = email.VerifyEmailFormat(req.Email); err != nil {
			// 邮箱格式不正确
			zap.L().Error("email VerifyEmailFormat failed", zap.Error(err))
			return err
		}
	}

	// 判断邮箱是否重复注册
	_, exist, err = userDao.ExistUserEmail(req.Email)
	if err != nil {
		// 数据库查询错误
		zap.L().Error("mysql ExistUserEmai failed", zap.Error(err))
		return err
	}
	if exist {
		// 邮箱存在
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return e.ErrorEmailExist
	}
	// 2. 生成UID
	userID := snowflake.GenID()
	// fmt.Println(userID)

	// 3. 构造一个user示例
	// 获取默认头像
	defaultAvatar := settings.Conf.AppConfig.PhotoPathConfig
	user := &models.User{
		UserID:   userID,
		UserName: req.UserName,
		Age:      req.Age,
		Email:    req.Email,
		Password: req.Password,
		Gender:   req.Gender,
		Avatar:   defaultAvatar.DefaultAvatar,
	}
	// 4. 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		zap.L().Error("user.SetPassword failed", zap.Error(err))
		return
	}
	// 5. 保存到数据库
	err = userDao.InsertUser(user)
	if err != nil {
		zap.L().Error("mysql InsertUser failed", zap.Error(err))
	}
	return err
}

// 登录业务
func (s *UserSrv) UserLogin(ctx context.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	userDao := mysql.NewUserDao(ctx)
	// user = &models.User{
	// 	UserName: req.Username,
	// 	Password: req.Password,
	// }
	// 1. 判断用户是否存在
	user, exist, err := userDao.CheckUserExist(req.UserName)
	if err != nil {
		// 数据库查询错误
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return nil, err
	}
	if !exist { // 如果查询不到，返回相应的错误
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return nil, e.ErrorUserNotExist
	}
	// 将用户登录输入的名称和密码信息传入model层的user中
	if err = user.CheckPassword(req.Password); err != nil {
		// 传递的是指针，能拿到数据库中注册时原本生成的UserID和Username
		zap.L().Error("mysql Login failed", zap.Error(err))
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		zap.L().Error("jwt GenToken failed", zap.Error(err))
	}
	user.Token = token
	// 保存到redis中
	if err = redisCache.StorgeUserIdToken(token, user.UserName); err != nil {
		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
	}
	resp = &types.UserLoginResp{
		UserID:   user.UserID,
		UserName: user.UserName,
		Token:    user.Token,
	}
	return
}

// 登录发送邮箱验证码业务
func (s *UserSrv) SendEmailCode(ctx context.Context, req *types.UserSendEmailCodeReq) (err error) {
	noticeDao := mysql.NewNoticeDao(ctx)
	// 检查是否在1分钟内发送过邮件， 如果有发送过，需等待后才能发送
	err = redisCache.EmailCodeExists(req.UserEmail)
	if err != nil {
		zap.L().Error("send email code from redis in 1 minutes")
		fmt.Printf("send email code from redis in 1 minutes")
		return err
	}

	// 获取六位数邮箱验证码
	code := email.GetConfirmCode()

	// 将其存储至Redis中，由于Redis为KV键值对存储所以需要定义前缀方便使用
	err = redisCache.StorgeEmailCode(req.UserEmail, code)
	if err != nil {
		zap.L().Error("redisCache.StorgeEmailCode failed, err:%\v", zap.Error(err))
		fmt.Printf("redisCache.StorgeEmailCode failed, err:%\v", err)
		return err
	}

	// 发送邮件，此处为方便起见没有处理返回值
	// address = settings.Conf.EmailConfig.VaildEmail + token
	var notice *models.Notice
	notice, err = noticeDao.GetNoticeById(req.OperationType)
	if err != nil {
		zap.L().Error("mysql GetNoticeById failed", zap.Error(err))
		return err
	}
	mailStr := notice.Text
	var builder strings.Builder
	builder.WriteString(mailStr)
	builder.Write([]byte(code))
	// mailTex := strings.Replace(mailStr, "Email", code, -1)
	// 第一个参数内容 第二个参数接收方 第三个标题 第四个发送方
	err = email.Send(builder.String(), req.UserEmail, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
	if err != nil {
		zap.L().Error("email Send failed", zap.Error(err))
		return err
	}
	// 设置每个邮箱发送邮件的时间 此处设置为1分钟，由于Redis为KV键值对存储所以需要定义前缀方便使用
	err = redisCache.SendEmailCode(req.UserEmail, code)
	if err != nil {
		zap.L().Error("redisCache.SendEmailCode failed, err:%\v", zap.Error(err))
		fmt.Printf("redisCache.SendEmailCode failed, err:%\v", err)
		return err
	}
	return nil
}

// 用户邮箱验证码登录业务
func (s *UserSrv) LoginEmail(ctx context.Context, req *types.UserLoginEmailReq) (user *models.User, err error) {
	userDao := mysql.NewUserDao(ctx)
	// 将用户邮箱登录输入的邮箱传入dao层mysql中进行判断是否存在，并取出
	// 从数据库中获取用户信息（用户名和id）
	user, err = userDao.GetUserByEmail(req.UserEmail)
	// 如果不存在则错误
	if err != nil {
		zap.L().Error(" mysql.NotExistUserEmail failed, err:%\v", zap.Error(err))
		fmt.Printf(" mysql.NotExistUserEmail failed, err:%\v", err)
		return nil, err
	}
	userInfo := &models.User{
		UserID:   user.UserID,
		Email:    req.UserEmail,
		UserName: user.UserName,
	}
	// 校验验证码正确性
	err = redisCache.CheckEmailCode(req.UserEmail, userInfo.Email, req.Code)
	if err != nil {
		zap.L().Error("redisCache.CheckEmailCode failed, err:%\v", zap.Error(err))
		fmt.Printf("redisCache.CheckEmailCode failed, err:%\v", err)
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(userInfo.UserID, userInfo.UserName)
	if err != nil {
		zap.L().Error("jwt GenToken failed", zap.Error(err))
	}
	user.Token = token
	// 保存到redis中
	if err = redisCache.StorgeUserIdToken(token, userInfo.Email); err != nil {
		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
	}
	return
}

// 用户信息修改
func (s *UserSrv) Update(ctx context.Context, req *types.UserUpdateReq) (err error) {
	u, _ := ctl.GetUserInfo(ctx)
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(u.UserId)
	// 1. 查询登录的用户ID，不用在数据库中查询是否存在，因为是已经登录过的
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}
	user = &models.User{
		UserID:   u.UserId,
		UserName: req.Username,
		Age:      req.Age,
		Email:    req.Email,
		Password: req.Password,
		Gender:   req.Gender,
		Avatar:   upload.AvatarURL() + req.Avatar,
	}
	err = userDao.UpdateUser(u.UserId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return err
	}
	return err
}

// 上传用户头像
func (s *UserSrv) UploadAvatar(ctx context.Context, file multipart.File, fileSize int64, req *types.UserAvatar) (resp interface{}, err error) {
	// u, _ := ctl.GetUserInfo(ctx)
	u, _ := ctl.GetLoginUserID(ctx)
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(u.UserId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return nil, err
	}
	// 保存到本地
	var path string
	if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal { // 兼容两种存储方式
		path, err = upload.UploadAvatarToLocalStatic(file, u.UserId, user.UserName)
	} else { //保存到七牛云oss
		path, err = upload.UploadToQiNiuAvatar(file, user.UserName, fileSize)
	}

	if err != nil {
		zap.L().Error("upload avatar failed", zap.Error(err))
		return nil, err
	}
	user.Avatar = path
	err = userDao.UpdateUser(u.UserId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return nil, err
	}
	// fmt.Println(user)
	resp = &types.UserAvatar{
		UserID:   u.UserId,
		UserName: user.UserName,
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

	return resp, nil
}

// 用户发送邮箱验证码绑定与解绑
func (s *UserSrv) SendEmail(ctx context.Context, req *types.UserSendEmailReq) (err error) {
	// u, _ := ctl.GetUserInfo(ctx)
	u, _ := ctl.GetLoginUserID(ctx)
	userDao := mysql.NewUserDao(ctx)
	noticeDao := mysql.NewNoticeDao(ctx)
	var user *models.User
	user, err = userDao.GetUserByID(u.UserId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}

	token, err := jwt.GenerateEmailToken(req.OperationType, u.UserId, req.Email, user.Password)
	if err != nil {
		zap.L().Error("jwt GenerateEmailToken failed", zap.Error(err))
		return err
	}
	// notice := new(models.Notice)
	var notice *models.Notice
	notice, err = noticeDao.GetNoticeById(req.OperationType)
	if err != nil {
		zap.L().Error("mysql GetNoticeById failed", zap.Error(err))
		return err
	}
	var address string
	address = settings.Conf.EmailConfig.VaildEmail + token
	mailStr := notice.Text
	mailTex := strings.Replace(mailStr, "Email", address, -1)
	err = email.Send(mailTex, req.Email, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
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

// 用户绑定邮箱和解绑邮箱业务
// TODO:传入参数优化
func (s *UserSrv) ValidEmail(ctx context.Context, token string) (err error) {
	var operationType int
	var uId int64
	var user *models.User
	var email string
	userDao := mysql.NewUserDao(ctx)
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
	user, err = userDao.GetUserByID(uId)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}
	if operationType == 1 {
		user.Email = email
	} else if operationType == 2 {
		user.Email = " "
	}
	err = userDao.UpdateUserEmail(uId, user)
	if err != nil {
		zap.L().Error("mysql UpdateUserEmail failed", zap.Error(err))
		return err
	}
	return
}
