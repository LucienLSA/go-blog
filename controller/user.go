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
// @Summary 用户注册接口
// @Description 根据用户所填信息进行注册
// @Tags 用户接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamSignUp true "用户注册"
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
// @Summary 用户登录接口
// @Description 根据所填用户名和密码进行登录
// @Tags 用户接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLogin true "用户登录"
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
		l := service.GetUserSrv()
		user, err := l.UserLogin(c.Request.Context(), &req)
		if err != nil {
			zap.L().Error("service login failed", zap.String("username", req.UserName), zap.Error(err))
			if errors.Is(err, e.ErrorUserNotExist) {
				e.ResponseError(c, e.CodeUserNotExist)
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
func SendEmailCodeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserSendEmailCodeReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
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
		err := l.SendEmailCode(ctx, &req)
		if err != nil {
			zap.L().Error("service SendEmailCode failed", zap.String("UserEmail", req.UserEmail), zap.Error(err))
			e.ResponseError(ctx, e.CodeServerBusy)
			return
		}
		e.ResponseSuccessData(ctx, nil)
	}
}

// LoginEmailHandler 用户邮箱验证码登录
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
		user, err := l.LoginEmail(ctx, &req)
		if err != nil {
			zap.L().Error("service login failed", zap.String("UserEmail", req.UserEmail), zap.Error(err))
			if errors.Is(err, e.ErrorEmailNotExit) {
				e.ResponseError(ctx, e.CodeEmailNotExist)
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
func UpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取参数和参数校验
		var req types.UserUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			// 请求参数有误 直接返回响应
			zap.L().Error("Updatewith invalid param", zap.Error(err))
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
		// claimsID, _ := request.GetLoginUserID(c)
		if err := l.Update(c, &req); err != nil {
			zap.L().Error("service update failed", zap.Error(err))
			// 理论上这个错误不会发生，因为是在登录情况下执行的，一定是存在的
			if errors.Is(err, e.ErrorUserNotExist) {
				e.ResponseError(c, e.CodeUserNotExist)
				return
			}
			e.ResponseError(c, e.CodeServerBusy)
			return
		}
		// 3. 返回响应
		e.ResponseSuccessData(c, nil)
	}
}

// 用户头像上传
func UploadAvatarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取上传的文件参数
		file, fileHeader, _ := c.Request.FormFile("file")
		if fileHeader == nil {
			e.ResponseError(c, e.CodeUploadFile)
			zap.L().Error("UpLoad file failed")
			return
		}
		fileSize := fileHeader.Size
		// 2. 用户校验
		var req types.UserAvatar
		if err := c.ShouldBind(&req); err == nil {
			// 获取登录用户的id
			// claimID, _ := request.GetLoginUserID(c)
			l := service.GetUserSrv()
			resp, err := l.UploadAvatar(c, file, fileSize, &req)
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
func SendEmailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserSendEmailReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
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
		if err := l.SendEmail(ctx, &req); err != nil {
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

// ValidEmailHandler用户验证邮箱
func ValidEmailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var p types.UserVaildEmail
		// p := new(models.ParamVaildEmail)
		if err := ctx.ShouldBind(&p); err != nil {
			// 2. 请求参数有误 直接返回响应
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
		} else {
			// 验证邮箱业务
			// TODO： 邮箱验证的token需要放在更隐蔽的未知
			token := ctx.Query("token")
			if token == " " {
				zap.L().Error("get token failed", zap.Error(err))
				e.ResponseError(ctx, e.CodeTokenInvalid)
				return
			}
			l := service.GetUserSrv()
			err := l.ValidEmail(ctx.Request.Context(), token)
			if err != nil {
				zap.L().Error("service VaildEmail failed", zap.Error(err))
				e.ResponseError(ctx, e.CodeServerBusy)
				return
			}
			// 3. 成功返回数据
			e.ResponseSuccessData(ctx, nil)
		}
	}
}
