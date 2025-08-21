package service

import (
	"context"
	"errors"
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
func (s *UserSrv) UserLogin(ctx context.Context, ip, ua string, req *types.UserLoginReq) (resp interface{}, err error) {
	// 登录限流逻辑暂时移除，如需启用请补齐 Redis 限流实现
	// 登录逻辑
	userDao := mysql.NewUserDao(ctx)
	user, exist, err := userDao.CheckUserExist(req.UserName)
	if err != nil {
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return nil, err
	}
	if !exist {
		// 登录失败
		zap.L().Error("userDao CheckUserExist failed", zap.Error(err))
		return nil, e.ErrorUserNotExist
	}

	// 单会话策略：允许新登录覆盖旧会话，不拦截已登录状态

	zap.L().Info("登录校验", zap.String("输入密码", req.Password), zap.String("数据库摘要", user.PasswordDigest))
	if err = user.CheckPassword(req.Password); err != nil {
		// 密码错误
		zap.L().Error("mysql Login failed", zap.String("输入密码", req.Password), zap.String("数据库摘要", user.PasswordDigest), zap.Error(err))
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		zap.L().Error("jwt GenToken failed", zap.Error(err))
	}
	user.Token = token
	// 保存到redis中
	if err = redisCache.StorgeUserToken(token, user.UserName); err != nil {
		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
	}
	// 登录成功
	resp = &types.UserLoginResp{
		UserID:   user.UserID,
		UserName: user.UserName,
		Token:    user.Token,
	}
	return
}

// UserLogout 用户登出
func (s *UserSrv) UserLogout(ctx context.Context, ip string, req *types.UserLogoutReq) (resp interface{}, err error) {
	// 从上下文获取当前登录用户，避免越权登出
	u, getErr := ctl.GetUserInfo(ctx)
	if getErr != nil || u == nil {
		zap.L().Error("GetUserInfo failed", zap.Error(getErr))
		return nil, e.ErrorNeedLogin
	}
	username := u.UserName
	// 从redis中删除当前用户token
	if err = redisCache.DeleteUserToken(username); err != nil {
		zap.L().Error("redisCache.DeleteUserToken failed", zap.Error(err))
		return nil, err
	}

	zap.L().Info("用户登出成功", zap.String("用户名", username), zap.String("IP", ip))
	return nil, nil
}

// 登录发送邮箱验证码业务
func (s *UserSrv) SendEmailCode(ctx context.Context, req *types.UserSendEmailCodeReq) (err error) {
	noticeDao := mysql.NewNoticeDao(ctx)
	// 检查是否在1分钟内发送过邮件， 如果有发送过，需等待后才能发送
	err = redisCache.EmailCodeExists(req.UserEmail, req.OperationType)
	if err != nil {
		zap.L().Error("send email code from redis in 1 minutes")
		fmt.Printf("send email code from redis in 1 minutes")
		return err
	}

	// 获取六位数邮箱验证码
	code := email.GetConfirmCode()

	// 将其存储至Redis中，由于Redis为KV键值对存储所以需要定义前缀方便使用
	err = redisCache.StorgeEmailCode(req.UserEmail, code, req.OperationType)
	if err != nil {
		zap.L().Error("redisCache.StorgeEmailCode failed", zap.Error(err))
		fmt.Printf("redisCache.StorgeEmailCode failed, err:%v", err)
		return err
	}

	// 发送邮件，此处为方便起见没有处理返回值
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

	// 发送邮件
	err = email.Send(builder.String(), req.UserEmail, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
	if err != nil {
		zap.L().Error("email Send failed", zap.Error(err))
		return err
	}
	// 设置每个邮箱发送邮件的时间 此处设置为1分钟，由于Redis为KV键值对存储所以需要定义前缀方便使用
	err = redisCache.SendEmailCode(req.UserEmail, code, req.OperationType)
	if err != nil {
		zap.L().Error("redisCache.SendEmailCode failed", zap.Error(err))
		fmt.Printf("redisCache.SendEmailCode failed, err:%v", err)
		return err
	}
	return nil
}

// 用户邮箱验证码登录业务
func (s *UserSrv) LoginEmail(ctx context.Context, ip string, req *types.UserLoginEmailReq) (user *models.User, err error) {
	userDao := mysql.NewUserDao(ctx)
	// 从数据库中获取用户信息（用户名和id）
	user, err = userDao.GetUserByEmail(req.UserEmail)
	// 如果不存在则错误
	if err != nil {
		zap.L().Error("mysql.NotExistUserEmail failed", zap.Error(err))
		fmt.Printf("mysql.NotExistUserEmail failed, err:%v", err)
		return nil, err
	}

	// 单会话策略：允许新登录覆盖旧会话，不拦截已登录状态

	userInfo := &models.User{
		UserID:   user.UserID,
		Email:    req.UserEmail,
		UserName: user.UserName,
	}

	// 校验验证码正确性
	err = redisCache.CheckEmailCode(req.UserEmail, req.Code, 3) // 3表示登录操作
	if err != nil {
		zap.L().Error("redisCache.CheckEmailCode failed", zap.Error(err))
		fmt.Printf("redisCache.CheckEmailCode failed, err:%v", err)
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(userInfo.UserID, userInfo.UserName)
	if err != nil {
		zap.L().Error("jwt GenToken failed", zap.Error(err))
	}
	user.Token = token
	// 保存到redis中
	if err = redisCache.StorgeUserToken(token, userInfo.UserName); err != nil {
		zap.L().Error("redisCache.StorgeUserIdToken failed", zap.Error(err))
	}
	return
}

// 用户信息修改
func (s *UserSrv) Update(ctx context.Context, req *types.UserUpdateReq) (err error) {
	// 获取当前登录用户信息
	u, err := ctl.GetUserInfo(ctx)
	if err != nil || u == nil {
		zap.L().Error("GetUserInfo failed", zap.Error(err))
		return errors.New("用户未登录或信息获取失败")
	}
	uid := u.UserId

	// 获取用户当前信息
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(uid)
	// 1. 查询登录的用户ID，不用在数据库中查询是否存在，因为是已经登录过的
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return err
	}

	// 如果提供了新密码，需要验证并更新
	if req.Password != "" && req.NewPassword != "" {
		if err = user.CheckPassword(req.Password); err != nil {
			return errors.New("当前密码错误")
		}
		if err = user.SetPassword(req.NewPassword); err != nil {
			return errors.New("新密码设置失败")
		}
		zap.L().Info("新密码加密摘要", zap.String("digest", user.PasswordDigest))
	}

	// 构建只包含需要更新字段的map
	updates := map[string]interface{}{}
	if req.Age > 0 {
		updates["age"] = req.Age
		user.Age = req.Age
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
		user.Gender = req.Gender
	}
	// 邮箱字段特殊处理：如果请求中包含email字段，则更新（包括空字符串，表示解绑）
	updates["email"] = req.Email
	user.Email = req.Email
	if req.Avatar != "" {
		updates["avatar"] = upload.AvatarURL() + req.Avatar
		user.Avatar = upload.AvatarURL() + req.Avatar
	}
	if req.Password != "" && req.NewPassword != "" {
		updates["password_digest"] = user.PasswordDigest
	}
	if req.UserName != "" {
		updates["user_name"] = req.UserName
		user.UserName = req.UserName
	}
	fmt.Println(user)
	fmt.Println(updates)
	zap.L().Info("用户信息更新字段", zap.Any("updates", updates))
	// 保存到数据库
	err = userDao.UpdateUser(uid, updates)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		return err
	}
	return nil
}

// 上传用户头像
func (s *UserSrv) UploadAvatar(ctx context.Context, file multipart.File, fileSize int64, fileName string, req *types.UserAvatar) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil || u == nil {
		zap.L().Error("GetUserInfo failed", zap.Error(err))
		return nil, errors.New("用户未登录或信息获取失败")
	}
	uid := u.UserId
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(uid)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return nil, errors.New("获取用户信息失败")
	}

	// 保存旧头像路径，用于失败时恢复
	oldAvatar := user.Avatar

	// 保存到本地
	var path string
	if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal { // 兼容两种存储方式
		path, err = upload.UploadAvatarToLocalStatic(file, uid, user.UserName, fileName)
	} else { //保存到七牛云oss
		path, err = upload.UploadToQiNiuAvatar(file, user.UserName, fileSize, fileName)
	}

	if err != nil {
		zap.L().Error("upload avatar failed", zap.Error(err))
		return nil, errors.New("头像上传失败")
	}

	// 更新数据库
	user.Avatar = path
	updates := map[string]interface{}{
		"avatar": upload.AvatarURL() + path,
	}
	err = userDao.UpdateUser(uid, updates)
	if err != nil {
		zap.L().Error("mysql UpdateUser failed", zap.Error(err))
		// 数据库更新失败，删除已上传的文件
		if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal {
			upload.DeleteOldAvatar(path)
		}
		return nil, errors.New("更新用户头像信息失败")
	}

	// 数据库更新成功，删除旧头像文件（如果是本地存储）
	if settings.Conf.AppConfig.UploadModel == settings.UploadModelLocal && oldAvatar != "" && oldAvatar != path {
		upload.DeleteOldAvatar(oldAvatar)
	}

	resp = &types.UserAvatar{
		UserID:   uid,
		UserName: user.UserName,
		Avatar:   upload.AvatarURL() + user.Avatar,
	}
	return resp, nil
}

// 用户发送邮箱验证码绑定与解绑
func (s *UserSrv) SendEmail(ctx context.Context, req *types.UserSendEmailReq) (err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil || u == nil {
		zap.L().Error("GetUserInfo failed", zap.Error(err))
		return errors.New("用户未登录或信息获取失败")
	}

	// 生成验证码并保存到 Redis
	code := email.GetConfirmCode()
	err = redisCache.StorgeEmailCode(req.Email, code, req.OperationType)
	if err != nil {
		zap.L().Error("redis StorgeEmailCode failed", zap.Error(err))
		return err
	}

	// 邮件正文为验证码
	var operationText string
	if req.OperationType == 1 {
		operationText = "绑定"
	} else {
		operationText = "解绑"
	}

	mailTex := fmt.Sprintf(`
	  <p>您正在%s邮箱, 验证码为: </p>
	  <h2 style=\"color:blue;\">%s</h2>
	  <p>请在页面输入该验证码完成操作。</p>
	`, operationText, code)
	if settings.Conf.AppConfig.Mode == "pro" {
		err = email.Send(mailTex, req.Email, settings.Conf.AppConfig.Name, settings.Conf.EmailConfig.SmtpEmail)
		if err != nil {
			zap.L().Error("email Send failed", zap.Error(err))
			return err
		}
	} else {
		return
	}
	return
}

// 用户绑定邮箱和解绑邮箱业务
// TODO:传入参数优化
func (s *UserSrv) ValidEmail(ctx context.Context, token string) (err error) {
	var operationType int
	var uId int64
	// var user *models.User // 已无用，删除
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
	// user, err = userDao.GetUserByID(uId) // 已无用，删除
	// if err != nil {
	// 	zap.L().Error("mysql GetUserByID failed", zap.Error(err))
	// 	return err
	// }
	updates := map[string]interface{}{}
	if operationType == 1 {
		updates["email"] = email
	} else if operationType == 2 {
		updates["email"] = " "
	}
	err = userDao.UpdateUser(uId, updates)
	if err != nil {
		zap.L().Error("mysql UpdateUserEmail failed", zap.Error(err))
		return err
	}
	return
}

// ValidEmailCode 验证邮箱验证码
func (s *UserSrv) ValidEmailCode(ctx context.Context, p *types.UserVaildEmail) error {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil || u == nil {
		zap.L().Error("GetUserInfo failed", zap.Error(err))
		return errors.New("用户未登录或信息获取失败")
	}
	uid := u.UserId
	// 获取用户信息
	userDao := mysql.NewUserDao(ctx)
	user, err := userDao.GetUserByID(uid)
	if err != nil {
		zap.L().Error("mysql GetUserByID failed", zap.Error(err))
		return errors.New("获取用户信息失败")
	}

	// 验证码校验
	if err := redisCache.CheckEmailCode(p.Email, p.Code, p.OperationType); err != nil {
		return err
	}

	// 如果是绑定操作，检查邮箱是否已被其他用户使用
	updates := map[string]interface{}{}
	if p.OperationType == 1 {
		// 绑定操作：检查邮箱是否已被其他用户使用
		existingUser, exist, err := userDao.ExistUserEmail(p.Email)
		if err != nil {
			zap.L().Error("mysql ExistUserEmail failed", zap.Error(err))
			return errors.New("检查邮箱状态失败")
		}
		if exist && existingUser.UserID != user.UserID {
			return errors.New("该邮箱已被其他用户绑定")
		}
		// 如果邮箱不存在或被当前用户绑定，允许绑定
		updates["email"] = p.Email
	} else if p.OperationType == 2 {
		// 解绑操作：验证用户是否拥有该邮箱
		if user.Email != p.Email {
			return errors.New("只能解绑自己的邮箱")
		}
		updates["email"] = ""
	} else {
		return errors.New("无效的操作类型")
	}

	// 更新用户信息
	return userDao.UpdateUser(user.UserID, updates)
}
