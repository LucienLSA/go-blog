package models

import (
	"blog-service/global"

	"gorm.io/gorm"
)

type Tag struct {
	Model
	CreatedBy  string `json:"created_by" gorm:"column:created_by; varchar(100)"`  // 创建人
	ModifiedBy string `json:"modified_by" gorm:"varchar(100);column:modified_by"` // 修改人
	Name       string `json:"name" gorm:"column:name; varchar(100)"`
	State      int    `json:"state" gorm:"column:state; tinyint(3)"`
}

// 表名
func (t Tag) TableName() string {
	return "blog_tag"
}

// 获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	if pageSize > 0 && pageNum > 0 {
		err = global.DBEngine.Where(maps).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = global.DBEngine.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

// 获取标签总数
func GetTagTotal(maps interface{}) (count int64, err error) {
	err = global.DBEngine.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 判断标签名重复
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := global.DBEngine.Select("id").Where("name =?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 新增标签
func CreateTags(name string, state int, createdBy string) error {
	global.DBEngine.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return nil
}

// 根据ID判断是否存在标签
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	// err := global.DBEngine.Select("id").Where("id =?", id).First(&tag).Error
	err := global.DBEngine.Select("id").Where("id =?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// logging.LogrusObj.Infoln(err)
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, err
}

// 更新标签
func UpdateTags(id int, data interface{}) error {
	err := global.DBEngine.Model(&Tag{}).Where("id =?", id).Updates(data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// logging.LogrusObj.Infoln(err)
		return err
	}
	return nil
}

// 删除标签 (软删除)
func DeleteTags(id int) error {
	err := global.DBEngine.Model(&Tag{}).Where("id =?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除标签（硬删除）
func CleanAllTags() bool {
	err := global.DBEngine.Unscoped().Delete(&Tag{}).Error
	if err != nil {
		return false
	}
	return true
}

// func (tag *Tag) BeforeCreateT(tx *gorm.DB) error {
// 	now := time.Now()
// 	tag.CreatedOn = &now
// 	return nil
// }

// func (tag *Tag) BeforeUpdateT(tx *gorm.DB) error {
// 	now := time.Now()
// 	tag.ModifiedOn = &now
// 	return nil
// }

// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//     scope.SetColumn("CreatedOn", time.Now().Unix())

//     return nil
// }

// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//     scope.SetColumn("ModifiedOn", time.Now().Unix())

//     return nil
// }
