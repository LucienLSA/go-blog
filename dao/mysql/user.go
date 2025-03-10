package mysql

import (
	"database/sql"

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
		return ErrorUserExist
	}
	return nil
}

// 检查指定用户名的用户是否不存在
func CheckUserNototExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	// fmt.Println(count)
	if count <= 0 {
		return ErrorUserNotExist
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
	sqlStr := `insert into user(user_id, username, password, email, age, gender, avatar) values(?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email, user.Age, user.Gender, user.Avatar)
	return err
}

// 向数据库更新用户信息
func UpdateUser(uId int64, user *models.User) (err error) {
	user.Password, err = SetPassword(user, user.Password)
	if err != nil {
		return err
	}
	sqlStr := `
	update user 
	set username=?,
		password=?,
		email=?,
		age=?,
		gender=?
	where user_id = ?`
	_, err = db.Exec(sqlStr, user.Username, user.Password, user.Email, user.Age, user.Gender, uId)
	return err
}

// 向数据库更新用户头像
func UpdateAvatar(uId int64, user *models.User) (err error) {
	userAvatar := user.Avatar
	sqlStr := `update user set avatar=? where user_id = ?`
	_, err = db.Exec(sqlStr, userAvatar, uId)
	return err
}

// 向数据库更新用户邮箱
func UpdateUserEmail(uId int64, user *models.User) (err error) {
	userEmail := user.Email
	sqlStr := `update user set email=? where user_id = ?`
	_, err = db.Exec(sqlStr, userEmail, uId)
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

// 检查指定用户名的邮箱是否存在
func ExistUserEmail(email string) (err error) {
	sqlStr := "select count(user_id) from user where email = ?"
	var count int
	err = db.Get(&count, sqlStr, email)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorEmailExist
	}
	return
}

// // 检查指定用户名的邮箱不存在，存在则获取用户
// func GetExistUserEmail(email string) (user *models.User, err error) {
// 	sqlStr := "select count(user_id) from user where email = ?"
// 	var count int
// 	err = db.Get(&count, sqlStr, email)
// 	if err != nil {
// 		return err
// 	}
// 	if count <= 0 {
// 		return ErrorEmailNotExit
// 	}
// 	return
// }

// // 邮箱登录 与mysql中的用户邮箱比对
// func LoginEmail(user *models.User) (err error) {
// 	// 记录用户输入的密码
// 	userEmail := user.Email
// 	// 执行SQL语句
// 	sqlStr := "select user_id from user where email=?"
// 	err = db.Get(user, sqlStr, user.Email)
// 	if err == sql.ErrNoRows { // sql自带查询错误
// 		return ErrorUserNotExist
// 	}
// 	if err != nil {
// 		// 查询数据库失败
// 		return err
// 	}
// 	// fmt.Println("用户输入密码", userPassword)
// 	// fmt.Println("数据库中密码", user.Password)
// 	// 判断邮箱是否正确
// 	return
// }

// 根据作者/用户id获取用户id和用户名
func GetUserByID(uId int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uId)
	return
}

// 根据作者/用户邮箱获取用户id和用户名
func GetUserByEmail(email string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where email = ?`
	err = db.Get(user, sqlStr, email)
	if err == sql.ErrNoRows { // sql自带查询错误
		return nil, ErrorEmailNotExit
	}
	if err != nil {
		// 查询数据库失败
		return nil, err
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
