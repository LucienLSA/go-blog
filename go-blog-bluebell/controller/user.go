package controller

import (
	"errors"

	"github.com/LucienLSA/go-blog/pkg/e"
	"github.com/LucienLSA/go-blog/pkg/email"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
	"github.com/LucienLSA/go-blog/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserSignUpReq true "注册信息"
// @Success 200 {object} _ResponseUserLogin
// @Router /user/signup [post]
// 注册请求
func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求参数和参数校验
		// var p models.ParamSignUp
		// p := new(models.ParamSignUp)
		var req types.UserSignUpReq
		if err := c.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("SignUp with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				// c.JSON(http.StatusOK, gin.H{
				// 	"msg": err.Error(),
				// })
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": translator.RemoveTopStruct(errs.Translate(translator.Trans)),
			// })
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// // 手动对请求参数进行详细业务规则的校验
		// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
		// 	zap.L().Error("SignUp with invalid param")
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"msg": "请求参数有误",
		// 	})
		// 	return
		// }
		// fmt.Println(p)
		// 2. 业务处理
		// 获取服务单例实例
		l := service.GetUserSrv()
		if err := l.UserSignUp(c.Request.Context(), &req); err != nil {
			zap.L().Error("service signup failed", zap.Error(err))
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": "注册失败",
			// 	// "error": err,
			// })
			if errors.Is(err, e.ErrorUserExist) {
				e.ResponseError(c, e.CodeUserExist)
				return
			}
			if errors.Is(err, e.ErrorEmailExist) {
				e.ResponseError(c, e.CodeEmailExist)
				return
			}
			if errors.Is(err, email.ErrorEmailFormat) {
				e.ResponseError(c, e.CodeEmailFormat)
				return
			}
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 3. 返回响应
		e.ResponseSuccessData(c, nil)
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "注册成功",
		// })
	}
}

// LoginHandler 用户登录
// @Summary 用户登录
// @Description 用户名密码登录
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserLoginReq true "登录信息"
// @Success 200 {object} _ResponseUserLogin
// @Router /user/login [post]
// 登录请求
func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserLoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("Login with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// 2. 业务处理
		ip := c.RemoteIP()
		l := service.GetUserSrv()
		user, err := l.UserLogin(c.Request.Context(), ip, c.Request.UserAgent(), &req)

		if err != nil {
			zap.L().Error("service login failed", zap.String("username", req.UserName), zap.Error(err))
			if errors.Is(err, e.ErrorUserNotExist) {
				e.ResponseError(c, e.CodeUserNotExist)
				return
			}
			if errors.Is(err, e.ErrorUserAlreadyLogin) {
				e.ResponseError(c, e.CodeUserAlreadyLogin)
				return
			}
			e.ResponseError(c, e.CodeInvalidPassword)
			return
		}
		// 3. 返回响应
		// e.ResponseSuccessData(c, gin.H{
		// 	"user_id":   user.ID, // id值大于1<<53-1, int64类型最大值为1<<63-1，转化为字符串
		// 	"user_name": user.UserName,
		// 	"token":     user.Token,
		// })
		e.ResponseSuccessData(c, user)
	}
}

// SendEmailCodeHandler 用户登录时给邮箱发送验证码
// @Summary 发送登录验证码
// @Description 发送用于登录的邮箱验证码
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserSendEmailCodeReq true "邮箱信息"
// @Success 200 {object} _ResponseSuccess
// @Router /user/email/login [post]
func SendEmailCodeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserSendEmailCodeReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("send email code with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(ctx, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(ctx, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// 2. 发送邮箱验证码业务
		l := service.GetUserSrv()
		err := l.SendEmailCode(ctx.Request.Context(), &req)
		if err != nil {
			zap.L().Error("service SendEmailCode failed", zap.String("UserEmail", req.UserEmail), zap.Error(err))
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(ctx, nil)
	}
}

// LoginEmailHandler 用户邮箱验证码登录
// @Summary 邮箱登录
// @Description 使用邮箱和验证码登录
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserLoginEmailReq true "邮箱登录信息"
// @Success 200 {object} _ResponseUserLogin
// @Router /user/email/login [post]
func LoginEmailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserLoginEmailReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("Login with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(ctx, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(ctx, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// 2. 业务处理
		l := service.GetUserSrv()
		ip := ctx.RemoteIP()
		user, err := l.LoginEmail(ctx.Request.Context(), ip, &req)
		if err != nil {
			zap.L().Error("service login failed", zap.String("UserEmail", req.UserEmail), zap.Error(err))
			if errors.Is(err, e.ErrorEmailNotExit) {
				e.ResponseError(ctx, e.CodeEmailNotExist)
				return
			}
			if errors.Is(err, e.ErrorUserAlreadyLogin) {
				e.ResponseError(ctx, e.CodeUserAlreadyLogin)
				return
			}
			e.ResponseError(ctx, e.CodeInvalidParam)
			return
		}
		e.ResponseSuccessData(ctx, gin.H{
			"user_id":   user.ID, // id值大于1<<53-1, int64类型最大值为1<<63-1，转化为字符串
			"user_name": user.UserName,
			"email":     req.UserEmail,
			"token":     user.Token,
		})
	}
}

// UpdateHandler 用户信息更新
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body types.UserUpdateReq true "更新信息"
// @Success 200 {object} _ResponseSuccess
// @Router /user/update [put]
func UpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("Update with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}

		// 2. 业务处理
		l := service.GetUserSrv()
		if err := l.Update(c.Request.Context(), &req); err != nil {
			zap.L().Error("service update failed", zap.Error(err))
			if errors.Is(err, e.ErrorUserNotExist) {
				e.ResponseError(c, e.CodeUserNotExist)
				return
			}
			if errors.Is(err, e.ErrorInvalidPassword) {
				e.ResponseError(c, e.CodeInvalidPassword)
				return
			}
			if errors.Is(err, e.ErrorEmailExist) {
				e.ResponseError(c, e.CodeEmailExist)
				return
			}
			if err.Error() == "用户未登录或信息获取失败" {
				e.ResponseError(c, e.CodeNeedLogin)
				return
			}
			e.ResponseError(c, e.CodeServerBusy)
			return
		}

		// 3. 返回响应
		e.ResponseSuccessData(c, nil)
	}
}

// UploadAvatarHandler 用户头像上传
// @Summary 上传头像
// @Description 为当前登录用户上传头像
// @Tags 用户接口
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param file formData file true "头像文件"
// @Success 200 {object} _ResponseSuccess
// @Router /user/avatar [post]
func UploadAvatarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取上传的文件参数
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil || fileHeader == nil {
			e.ResponseError(c, e.CodeUploadFile)
			zap.L().Error("UpLoad file failed", zap.Error(err))
			return
		}

		// 2. 文件验证
		fileSize := fileHeader.Size
		fileName := fileHeader.Filename

		// 验证文件大小 (限制为5MB)
		const maxFileSize = 5 * 1024 * 1024 // 5MB
		if fileSize > maxFileSize {
			e.ResponseError(c, e.CodeFileTooLarge)
			zap.L().Error("File size too large", zap.Int64("size", fileSize), zap.String("filename", fileName))
			return
		}

		// 验证文件类型
		contentType := fileHeader.Header.Get("Content-Type")
		allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
		isValidType := false
		for _, allowedType := range allowedTypes {
			if contentType == allowedType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			e.ResponseError(c, e.CodeInvalidFileType)
			zap.L().Error("Invalid file type", zap.String("contentType", contentType), zap.String("filename", fileName))
			return
		}

		// 3. 用户校验
		var req types.UserAvatar
		if err := c.ShouldBind(&req); err == nil {
			// 获取登录用户的id
			// claimID, _ := request.GetLoginUserID(c)
			l := service.GetUserSrv()
			resp, err := l.UploadAvatar(c.Request.Context(), file, fileSize, fileName, &req)
			if err != nil {
				zap.L().Error("service UploadAvatar failed, err:", zap.Error(err))
				e.ResponseError(c, e.CodeServerBusy)
				return
			}
			e.ResponseSuccessData(c, resp)
		} else { // 3.请求参数有误 直接返回响应
			zap.L().Error("Login with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
	}
}

// SendEmailHandler 发送邮箱验证码及信息
// @Summary 发送绑定/解绑邮箱验证码
// @Description 发送用于绑定或解绑邮箱的验证码
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body types.UserSendEmailReq true "操作信息"
// @Success 200 {object} _ResponseSuccess
// @Router /user/email/send [post]
func SendEmailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserSendEmailReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("SendEmail with invalid param", zap.Error(err))
			// 判断err是不是validator.ValidationErrors类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(ctx, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(ctx, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// claimsID, _ := request.GetLoginUserID(ctx)
		l := service.GetUserSrv()
		if err := l.SendEmail(ctx.Request.Context(), &req); err != nil {
			zap.L().Error("service SendEmail failed", zap.Error(err))
			if errors.Is(err, e.ErrorUserNotExist) {
				e.ResponseError(ctx, e.CodeUserNotExist)
				return
			}
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		// 3. 返回响应
		e.ResponseSuccessData(ctx, nil)
	}
}

// ValidEmailHandler 用户验证邮箱(绑定或解绑)
// @Summary 用户验证邮箱
// @Description 用户验证邮箱(绑定或解绑)
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserVaildEmail true "验证邮箱信息"
// @Success 200 {object} _ResponseUserLogin
// @Router /user/valid_email [post]
func ValidEmailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求参数和参数校验
		var req types.UserVaildEmail
		if err := c.ShouldBind(&req); err != nil {
			zap.L().Error("ValidEmail with invalid param", zap.Error(err))
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				e.ResponseError(c, e.CodeInvalidParam)
				return
			}
			e.ResponseErrorMsg(c, e.CodeInvalidParam,
				translator.RemoveTopStruct(errs.Translate(translator.Trans)))
			return
		}
		// 2. 业务处理
		l := service.GetUserSrv()
		if err := l.ValidEmailCode(c.Request.Context(), &req); err != nil {
			zap.L().Error("service valid email failed", zap.Error(err))
			if errors.Is(err, e.ErrorEmailExist) {
				e.ResponseError(c, e.CodeEmailExist)
				return
			}
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 3. 返回响应
		e.ResponseSuccessData(c, nil)
	}
}

// LogoutHandler 用户登出
// @Summary 用户登出
// @Description 用户登出，清除token
// @Tags 用户接口
// @Accept json
// @Produce json
// @Param data body types.UserLogoutReq true "登出信息"
// @Success 200 {object} _ResponseUserLogin
// @Router /user/logout [post]
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 2. 获取客户端IP
		ip := c.ClientIP()

		// 3. 业务处理
		l := service.GetUserSrv()
		if _, err := l.UserLogout(c.Request.Context(), ip, &types.UserLogoutReq{}); err != nil {
			zap.L().Error("service logout failed", zap.Error(err))
			e.ResponseError(c, e.CodeServerBusy)
			return
		}

		// 4. 返回响应
		e.ResponseSuccessData(c, nil)
	}
}
