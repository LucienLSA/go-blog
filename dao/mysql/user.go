package mysql

import (
	"database/sql"
	"errors"

	"github.com/LucienLSA/go-blog/models"
	"golang.org/x/crypto/bcrypt"
)

// 每一步数据库操作封装成函数
// 待service层业务需求进行调用

const PasswordCost = 12

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
)

// 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	// fmt.Println(count)
	if count > 0 {
		return ErrorUserExist
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

// 用户登录 与数据库中用户信息比对
func Login(user *models.User) (err error) {
	// 记录用户输入的密码
	userPassword := user.Password
	// 执行SQL语句
	sqlStr := "select user_id, username, password from user where username=?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows { // sql自带查询错误
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// fmt.Println("用户输入密码", userPassword)
	// fmt.Println("数据库中密码", user.Password)
	// 判断密码是否正确
	err = CheckPassword(user, userPassword)
	if err != nil {
		return ErrorInvalidPassword
	}
	return
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
	// 将原密码进行加密
	user.PasswordDigest = string(bytes)
	return user.PasswordDigest, nil
}

// CheckPassword 校验密码
func CheckPassword(user *models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
