package ctl

import (
	"context"
	"errors"
)

type key int

var userKey key

type UserInfo struct {
	UserId int64 `json:"user_id"`
}

var (
	ErrorUserNotLogin = errors.New("用户未登录")
	ErrorGetUserInfo  = errors.New("获取用户信息错误")
)

// // 获取到经过中间件的jwt auth鉴权后登录的用户信息 ，进行下一步处理
// func GetLoginUserID(ctx *gin.Context) (userID int64, err error) {
// 	uid, ok := ctx.Get(middlewares.CtxtUserIDKey)
// 	// u, ok := ctx.Value(middlewares.CtxtUserIDKey).(*UserInfo)
// 	if !ok {
// 		err = ErrorUserNotLogin
// 		return
// 	}
// 	userID, ok = uid.(int64)
// 	if !ok {
// 		err = ErrorUserNotLogin
// 		return
// 	}
// 	return
// }

// 从context获取用户信息
func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userKey).(*UserInfo)
	return u, ok
}

// 从context中获取用户信息id
func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, ErrorGetUserInfo
	}
	return user, nil
}

// 将用户信息存入context中
func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func InitUserInfo(ctx context.Context) {
	// TOOD 放缓存，之后的用户信息，走缓存
}
