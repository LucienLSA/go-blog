package mysql

import (
	"context"

	"github.com/LucienLSA/go-blog/repository/db/models"
	"gorm.io/gorm"
)

type CommunityDao struct {
	*gorm.DB
}

func NewCommunityDao(ctx context.Context) *CommunityDao {
	return &CommunityDao{NewDBClient(ctx)}
}

func NewCommunityDaoByDB(db *gorm.DB) *CommunityDao {
	return &CommunityDao{db}
}

func (dao *CommunityDao) GetCommunityList() (communityList []*models.Community, err error) {
	// sqlStr := "select community_id, community_name from community"
	// if err := db.Select(&communityList, sqlStr); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		zap.L().Warn("there is no community in database")
	// 		err = nil
	// 	}
	// }
	err = dao.DB.Model(&models.Community{}).Find(&communityList).Error
	return
}

func (dao *CommunityDao) GetCommunityDetailList(cid int64) (community *models.Community, err error) {
	community = new(models.Community)
	err = dao.DB.Model(&models.Community{}).Where("community_id=?", cid).First(community).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return community, nil
}

func (dao *CommunityDao) CreateCommunity(community *models.Community) (err error) {
	// sqlStr := "select community_id, community_name from community"
	// if err := db.Select(&communityList, sqlStr); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		zap.L().Warn("there is no community in database")
	// 		err = nil
	// 	}
	// }
	return dao.DB.Model(&models.Community{}).Create(&community).Error
}
