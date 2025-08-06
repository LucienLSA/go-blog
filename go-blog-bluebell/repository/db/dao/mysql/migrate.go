package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
)

func migrate() (err error) {
	err = _db.AutoMigrate(&models.User{}, &models.Post{}, &models.Notice{}, &models.Community{})
	if err != nil {
		return
	}

	// 迁移审核相关表
	err = MigrateReviewTables(context.Background())
	if err != nil {
		return
	}

	// 为帖子表添加审核字段
	err = AddReviewFieldsToPost(context.Background())
	if err != nil {
		return
	}

	// 初始化审核数据
	err = InitReviewData(context.Background())
	if err != nil {
		return
	}

	return nil
}
