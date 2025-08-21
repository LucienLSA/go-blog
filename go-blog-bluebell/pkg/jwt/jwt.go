package jwt

import (
	"errors"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"github.com/golang-jwt/jwt"
)

// const (
// 	AccessTokenExpireDuration  = time.Hour * 1
// 	RefreshTokenExpireDuration = time.Second * 30
// )

var MySecret = []byte("lucien-goblog-bluebell-jwt")

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID             int64  `json:"user_id"`
	Username           string `json:"username"`
	jwt.StandardClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(settings.Conf.AppConfig.JwtExpireTime * time.Minute)
	// 创建一个自己的声明数据
	claims := MyClaims{
		UserID:   userID,
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "bluebell", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	token, err := tokenClaims.SignedString(MySecret)

	return token, err
}

// func GenToken(userID int64, username string) (aToken, rToken string, err error) {
// 	nowTime := time.Now()
// 	aExpireTime := nowTime.Add(AccessTokenExpireDuration)
// 	rExpireTime := nowTime.Add(RefreshTokenExpireDuration)
// 	// 创建一个自己的声明数据
// 	claims := MyClaims{
// 		UserID:   userID,
// 		Username: username, // 自定义字段
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: aExpireTime.Unix(),
// 			Issuer:    "bluebell", // 签发人
// 		},
// 	}
// 	// 使用指定的签名方法创建签名对象
// 	atokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	// 使用指定的secret签名并获得完整的编码后的字符串token
// 	aToken, err = atokenClaims.SignedString(MySecret)

// 	// refresh token 不需要任何自定义数据
// 	rTokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
// 		ExpiresAt: rExpireTime.Unix(),
// 		Issuer:    "bluebell",
// 	})
// 	rToken, err = rTokenClaims.SignedString(MySecret)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return aToken, rToken, err
// }

// ParseToken 解析JWT
func ParseToken(tokenString string) (myclaims *MyClaims, err error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	myclaims = new(MyClaims) // 必须手动初始化 分配内存空间
	tokenClaims, err := jwt.ParseWithClaims(tokenString, myclaims, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if tokenClaims != nil {
		// 对token对象中的Claim进行类型断言
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid { // 校验toke
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}

// // RefreshToken 刷新Access Token
// func ParseRefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
// 	// access token 无效直接返回
// 	accessClaim, err := ParseToken(aToken)
// 	if err != nil {
// 		return
// 	}
// 	// refresh token 无效直接返回
// 	refreshClaim, err := ParseToken(rToken)
// 	if err != nil {
// 		return
// 	}

// 	if accessClaim.ExpiresAt > time.Now().Unix() {
// 		// 如果 access_token 没过期,每一次请求都刷新 refresh_token 和 access_token
// 		return GenToken(accessClaim.UserID, accessClaim.Username)
// 	}

// 	if refreshClaim.ExpiresAt > time.Now().Unix() {
// 		// 如果 access_token 过期了,但是 refresh_token 没过期, 刷新 refresh_token 和 access_token
// 		return GenToken(accessClaim.UserID, accessClaim.Username)
// 	}
// 	return "", "", errors.New("身份过期请重新登录")
// }
