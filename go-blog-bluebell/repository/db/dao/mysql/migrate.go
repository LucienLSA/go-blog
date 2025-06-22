package mysql

import "github.com/LucienLSA/go-blog/repository/db/models"

func migrate() (err error) {
	err = _db.AutoMigrate(&models.User{}, &models.Post{}, &models.Notice{}, &models.Community{})
	if err != nil {
		return
	}
	return nil
}
