package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// 每一步数据库操作封装成函数
// 待service层业务需求进行调用

// 检查指定用户名的用户是否存在
func (dao *UserDao) CheckUserExist(username string) (user *models.User, exist bool, err error) {
	// sqlStr := "select count(user_id) from user where username = ?"
	var count int64
	err = dao.DB.Model(&models.User{}).Where("user_name = ?", username).Count(&count).Error
	if err != nil {
		return user, false, err
	}
	if count <= 0 {
		return user, false, err
	}
	// if err = db.Get(&count, sqlStr, username); err != nil {
	// 	return err
	// }
	// fmt.Println(count)
	err = dao.DB.Model(&models.User{}).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// 检查指定用户名的邮箱是否存在
func (dao *UserDao) ExistUserEmail(email string) (user *models.User, exist bool, err error) {
	// sqlStr := "select count(user_id) from user where email = ?"
	var count int64
	// err = db.Get(&count, sqlStr, email)
	err = dao.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return user, false, err
	}
	if count <= 0 {
		return nil, false, err
	}
	err = dao.DB.Model(&models.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// 向数据库中插入一条用户信息
func (dao *UserDao) InsertUser(user *models.User) (err error) {
	// 执行SQL语句
	// sqlStr := `insert into user(user_id, username, password, email, age, gender, avatar) values(?,?,?,?,?,?,?)`
	// _, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email, user.Age, user.Gender, user.Avatar)
	// return err
	return dao.DB.Model(&models.User{}).Create(&user).Error
}

// // 检查指定用户名的用户是否不存在
// func CheckUserNototExist(username string) (err error) {
// 	sqlStr := "select count(user_id) from user where username = ?"
// 	var count int
// 	if err = db.Get(&count, sqlStr, username); err != nil {
// 		return err
// 	}
// 	// fmt.Println(count)
// 	if count <= 0 {
// 		return ErrorUserNotExist
// 	}
// 	return nil
// }

// 向数据库更新用户信息
func (dao *UserDao) UpdateUser(uId int64, user *models.User) (err error) {
	return dao.DB.Model(&models.User{}).Where("user_id=?", uId).
		Updates(&user).Error
	// user.Password, err = SetPassword(user, user.Password)
	// if err != nil {
	// 	return err
	// }
	// sqlStr := `
	// update user
	// set username=?,
	// 	password=?,
	// 	email=?,
	// 	age=?,
	// 	gender=?
	// where user_id = ?`
	// _, err = db.Exec(sqlStr, user.Username, user.Password, user.Email, user.Age, user.Gender, uId)
}

// // 向数据库更新用户头像
// func UpdateAvatar(uId int64, user *models.User) (err error) {
// 	userAvatar := user.Avatar
// 	sqlStr := `update user set avatar=? where user_id = ?`
// 	_, err = db.Exec(sqlStr, userAvatar, uId)
// 	return err
// }

// 向数据库更新用户邮箱
func (dao *UserDao) UpdateUserEmail(uId int64, user *models.User) (err error) {
	// userEmail := user.Email
	// sqlStr := `update user set email=? where user_id = ?`
	// _, err = db.Exec(sqlStr, userEmail, uId)
	return dao.DB.Model(&models.User{}).Where("user_id=?", uId).
		Updates(&user).Error
}

// // // 检查指定用户名的邮箱不存在，存在则获取用户
// // func GetExistUserEmail(email string) (user *models.User, err error) {
// // 	sqlStr := "select count(user_id) from user where email = ?"
// // 	var count int
// // 	err = db.Get(&count, sqlStr, email)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	if count <= 0 {
// // 		return ErrorEmailNotExit
// // 	}
// // 	return
// // }

// // // 邮箱登录 与mysql中的用户邮箱比对
// // func LoginEmail(user *models.User) (err error) {
// // 	// 记录用户输入的密码
// // 	userEmail := user.Email
// // 	// 执行SQL语句
// // 	sqlStr := "select user_id from user where email=?"
// // 	err = db.Get(user, sqlStr, user.Email)
// // 	if err == sql.ErrNoRows { // sql自带查询错误
// // 		return ErrorUserNotExist
// // 	}
// // 	if err != nil {
// // 		// 查询数据库失败
// // 		return err
// // 	}
// // 	// fmt.Println("用户输入密码", userPassword)
// // 	// fmt.Println("数据库中密码", user.Password)
// // 	// 判断邮箱是否正确
// // 	return
// // }

// 根据作者/用户id获取用户id和用户名
func (dao *UserDao) GetUserByID(uId int64) (user *models.User, err error) {
	// user = new(models.User)
	// sqlStr := `select user_id, username from user where user_id = ?`
	// err = db.Get(user, sqlStr, uId)
	err = dao.DB.Model(&models.User{}).Where("user_id=?", uId).
		First(&user).Error
	return

}

// 根据作者/用户邮箱获取用户id和用户名
func (dao *UserDao) GetUserByEmail(email string) (user *models.User, err error) {
	// sqlStr := `select user_id, username from user where email = ?`
	// err = db.Get(user, sqlStr, email)
	// if err == sql.ErrNoRows { // sql自带查询错误
	// 	return nil, ErrorEmailNotExit
	// }
	// if err != nil {
	// 	// 查询数据库失败
	// 	return nil, err
	// }
	err = dao.DB.Model(&models.User{}).Where("email=?", email).
		First(&user).Error
	return
}
