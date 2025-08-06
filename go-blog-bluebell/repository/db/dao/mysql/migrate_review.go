package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"go.uber.org/zap"
)

// MigrateReviewTables 迁移审核相关表
func MigrateReviewTables(ctx context.Context) error {
	db := NewDBClient(ctx)

	// 自动迁移表结构
	err := db.AutoMigrate(
		&models.PostReviewTask{},
		&models.PostReviewLog{},
	)

	if err != nil {
		zap.L().Error("migrate review tables failed", zap.Error(err))
		return err
	}

	zap.L().Info("migrate review tables success")
	return nil
}

// AddReviewFieldsToPost 为帖子表添加审核字段
func AddReviewFieldsToPost(ctx context.Context) error {
	db := NewDBClient(ctx)

	// 检查字段是否已存在
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'post' AND COLUMN_NAME = 'review_status'").Scan(&count).Error
	if err != nil {
		zap.L().Error("check review_status column failed", zap.Error(err))
		return err
	}

	// 如果字段不存在，则添加
	if count == 0 {
		// 添加审核相关字段
		sqlStatements := []string{
			"ALTER TABLE post ADD COLUMN review_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' COMMENT '审核状态'",
			"ALTER TABLE post ADD COLUMN review_result TEXT COMMENT '审核结果详情'",
			"ALTER TABLE post ADD COLUMN reviewed_at TIMESTAMP NULL COMMENT '审核时间'",
			"ALTER TABLE post ADD COLUMN review_score DECIMAL(3,2) DEFAULT 0.00 COMMENT '审核评分'",
			"ALTER TABLE post ADD COLUMN review_reason VARCHAR(500) COMMENT '拒绝原因'",
			"ALTER TABLE post ADD COLUMN review_suggestions TEXT COMMENT '修改建议'",
			"ALTER TABLE post ADD COLUMN review_tags JSON COMMENT '内容标签'",
		}

		for _, sql := range sqlStatements {
			err := db.Exec(sql).Error
			if err != nil {
				zap.L().Error("execute sql failed", zap.String("sql", sql), zap.Error(err))
				return err
			}
		}

		// 创建索引
		indexStatements := []string{
			"CREATE INDEX idx_review_status ON post (review_status)",
			"CREATE INDEX idx_reviewed_at ON post (reviewed_at)",
		}

		for _, sql := range indexStatements {
			err := db.Exec(sql).Error
			if err != nil {
				zap.L().Warn("create index failed", zap.String("sql", sql), zap.Error(err))
				// 索引创建失败不影响主要功能，只记录警告
			}
		}

		zap.L().Info("add review fields to post table success")
	} else {
		zap.L().Info("review fields already exist in post table")
	}

	return nil
}

// InitReviewData 初始化审核相关数据
func InitReviewData(ctx context.Context) error {
	db := NewDBClient(ctx)

	// 为现有帖子设置默认审核状态
	err := db.Exec("UPDATE post SET review_status = 'approved' WHERE review_status IS NULL OR review_status = ''").Error
	if err != nil {
		zap.L().Error("init review data failed", zap.Error(err))
		return err
	}

	zap.L().Info("init review data success")
	return nil
}
