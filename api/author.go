package api

import (
	"blog-service/pkg/app"
	"blog-service/pkg/e"
	"blog-service/pkg/util"
	"blog-service/service/auth_service"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

type Author struct {
	Username string `valid:"Required;MaxSize(20)"`
	Password string `valid:"Required;MaxSize(20)"`
}

func NewAuthor() *Author {
	return &Author{}
}

// @Summary 用户登录
// @Produce  json
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}} "
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}} "
// @Router /author/login [get]
func (au *Author) LoginAuthor(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	// valid.Required(username, "username").Message("用户名不能为空")
	// valid.Required(password, "password").Message("密码不能为空")
	// valid.MaxSize(username, 20, "username").Message("用户名长度不能超过20")
	// valid.MaxSize(password, 20, "password").Message("密码长度不能超过20")

	a := Author{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, valid.Errors)
		return
	}
	authService := auth_service.Author{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"token": token,
	})
	// data := make(map[string]interface{})
	// code := e.INVALID_PARAMS

	// if ok {
	// 	isExist := models.CheckAuthor(username, password)
	// 	if isExist {
	// 		token, err := util.GenerateToken(username, password)
	// 		if err != nil {
	// 			code = e.ERROR_AUTH_TOKEN
	// 			logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 		} else {
	// 			data["token"] = token
	// 			code = e.SUCCESS
	// 		}
	// 	} else {
	// 		code = e.ERROR_AUTH
	// 		logging.LogrusObj.Infoln(e.GetMsg(code)) // 补充错误处理
	// 	}
	// } else {
	// 	for _, err := range valid.Errors {
	// 		logging.LogrusObj.Infoln(err) // 补充错误处理
	// 		// log.Fatalf("err.key: %s, err.message: %s", err.Key, err.Message)
	// 	}
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  e.GetMsg(code),
	// 	"data": data,
	// })

}
