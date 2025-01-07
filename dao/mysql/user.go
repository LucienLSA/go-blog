package mysql

import (
	"errors"

	"github.com/LucienLSA/go-blog/models"
	"golang.org/x/crypto/bcrypt"
)

// 每一步数据库操作封装成函数
// 待service层业务需求进行调用

const PasswordCost = 12

// 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	// fmt.Println(count)
	if count > 0 {
		return errors.New("用户已存在")
	}
	return nil
}

// 向数据库中插入一条用户信息
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password, err = SetPassword(user, user.Password)
	if err != nil {
		return err
	}
	// fmt.Println(user.Password)
	// 执行SQL语句
	sqlStr := `insert into user(user_id, username, password, email, age, gender) values(?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email, user.Age, user.Gender)
	return err
}

// // 密码加密 旧版
// func encryptPassword(oPassword string) string {
// 	h := md5.New()
// 	h.Write([]byte(secret))
// 	return hex.EncodeToString(h.Sum([]byte(oPassword)))
// }

// SetPassword 设置密码
func SetPassword(user *models.User, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return "", err
	}
	user.PasswordDigest = string(bytes)
	return user.PasswordDigest, nil
}

// CheckPassword 校验密码
func CheckPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
