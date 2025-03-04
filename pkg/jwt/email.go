package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type EmailClaims struct {
	UserID        int64  `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType int    `json:"operation_type"`
	jwt.StandardClaims
}

// 生成邮箱验证token
func GenerateEmailToken(operation_type int, userID int64, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	claims := EmailClaims{
		UserID:        userID,
		Email:         email,
		Password:      password,
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
