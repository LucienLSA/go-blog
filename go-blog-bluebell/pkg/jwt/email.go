package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type EmailClaims struct {
	OperationType int    `json:"operation_type" form:"operation_type"`
	UserID        int64  `json:"user_id" form:"user_id"`
	Email         string `json:"email" form:"email"`
	jwt.StandardClaims
}

// 生成邮箱验证token
func GenerateEmailToken(operation_type int, userID int64, email string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	claims := EmailClaims{
		UserID:        userID,
		Email:         email,
		OperationType: operation_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "bluebell",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(MySecret)
	return token, err
}

// 验证邮件token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenEmailClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if tokenEmailClaims != nil {
		if claims, ok := tokenEmailClaims.Claims.(*EmailClaims); ok && tokenEmailClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
