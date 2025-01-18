package controller

import (
	"errors"
	"strconv"

	"github.com/LucienLSA/go-blog/dao/mysql"
	"github.com/LucienLSA/go-blog/models"
	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/LucienLSA/go-blog/pkg/translator"
	"github.com/LucienLSA/go-blog/service"
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
// @Router /signup [post]
// 注册请求
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	// var p models.ParamSignUp
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误 直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			response.ResponseError(c, response.CodeInvalidParam)
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": translator.RemoveTopStruct(errs.Translate(translator.Trans)),
		// })
		response.ResponseErrorMsg(c, response.CodeInvalidParam,
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
	if err := service.SignUp(p); err != nil {
		zap.L().Error("service signup failed", zap.Error(err))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "注册失败",
		// 	// "error": err,
		// })
		if errors.Is(err, mysql.ErrorUserExist) {
			response.ResponseError(c, response.CodeUserExist)
			return
		}
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.ResponseSuccessData(c, nil)
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "注册成功",
	// })
}

// LoginHandler 用户登录
// @Summary 用户登录接口
// @Description 根据所填用户名和密码进行登录
// @Tags 用户接口
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLogin true "用户登录"
// @Success 200 {object} _ResponseUserLogin
// @Router /login [post]
// 登录请求
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误 直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			response.ResponseError(c, response.CodeInvalidParam)
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": translator.RemoveTopStruct(errs.Translate(translator.Trans)),
		// })
		response.ResponseErrorMsg(c, response.CodeInvalidParam,
			translator.RemoveTopStruct(errs.Translate(translator.Trans)))
		return
	}
	// 2. 业务处理
	user, err := service.Login(p)
	if err != nil {
		zap.L().Error("service login failed", zap.String("username", p.Username), zap.Error(err))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "用户名或密码错误",
		// 	// "error": err,
		// })
		if errors.Is(err, mysql.ErrorUserNotExist) {
			response.ResponseError(c, response.CodeUserNotExist)
			return
		}
		response.ResponseError(c, response.CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登录成功",
	// })
	response.ResponseSuccessData(c, gin.H{
		"user_id":   strconv.FormatInt(user.UserID, 10), // id值大于1<<53-1, int64类型最大值为1<<63-1，转化为字符串
		"user_name": user.Username,
		"token":     user.Token,
	})
}
