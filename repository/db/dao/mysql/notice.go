package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"gorm.io/gorm"
)

type NoticeDao struct {
	*gorm.DB
}

func NewNoticeDao(ctx context.Context) *NoticeDao {
	return &NoticeDao{NewDBClient(ctx)}
}

func NewNoticeDaoByDB(db *gorm.DB) *NoticeDao {
	return &NoticeDao{db}
}

// 根据id获取Notice
func (dao *NoticeDao) GetNoticeById(operation_type int) (notice *models.Notice, err error) {
	// sqlStr := `select text from notice where id = ?`
	// err = db.Get(notice, sqlStr, operation_type)
	err = dao.DB.Model(&models.Notice{}).Where("id=?", operation_type).
		First(&notice).Error
	return
}
