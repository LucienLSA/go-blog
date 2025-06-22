package request

// var ErrorUserNotLogin = errors.New("用户未登录")

// // 获取到经过中间件的jwt auth鉴权后登录的用户信息 ，进行下一步处理
// func GetLoginUserID(ctx *gin.Context) (userID int64, err error) {
// 	uid, ok := ctx.Get(middlewares.CtxtUserIDKey)
// 	if !ok {
// 		err = ErrorUserNotLogin
// 		return
// 	}
// 	if userID != uid {
// 		err = ErrorUserNotLogin
// 		return
// 	}
// 	return
// }
