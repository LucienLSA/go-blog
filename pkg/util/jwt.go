package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// var jwtSecret = []byte(global.AppSetting.JwtSecret) // ? 为什么读不到

var jwtSecret = []byte("Lucien-go-blog-secret-jwt")

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 签发用户Token
func GenerateToken(password, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // 过期时间
	claims := Claims{
		Password: password,
		Username: username,
		// 声明标准JWT
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Lucien",
		},
	}
	// 签名并生成token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// 验证用户token
func ParseToken(token string) (*Claims, error) {
	// 解析Token 为Claims结构
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// 传入密钥
		return jwtSecret, nil
	})
	// 判断token是否有效
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
