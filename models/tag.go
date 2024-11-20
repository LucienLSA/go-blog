package models

import (
	"blog-service/global"
	"time"

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
	err = global.DBEngine.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil {
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
func ExitTagByName(name string) bool {
	var tag Tag
	global.DBEngine.Select("id").Where("name =?", name).First(&tag)
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return false
	// }
	if tag.ID > 0 {
		return true
	}
	return false
}

// 新增标签
func CreateTags(name string, state int, createdBy string) bool {
	global.DBEngine.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return true
}

// 根据ID判断是否存在标签
func ExistTagByID(id int) bool {
	var tag Tag
	err := global.DBEngine.Select("id").Where("id =?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tag.ID > 0 {
		return true
	}
	return false
}

// 更新标签
func UpdateTags(id int, data interface{}) bool {
	err := global.DBEngine.Model(&Tag{}).Where("id =?", id).Updates(data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	return true
}

// 删除标签
func DeleteTags(id int) bool {
	err := global.DBEngine.Model(&Tag{}).Where("id =?", id).Delete(&Tag{})
	if err.Error != nil {
		return false
	}
	return true
}

func (tag *Tag) BeforeCreateT(tx *gorm.DB) error {
	now := time.Now()
	tag.CreatedOn = &now
	return nil
}

func (tag *Tag) BeforeUpdateT(tx *gorm.DB) error {
	now := time.Now()
	tag.ModifiedOn = &now
	return nil
}

// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//     scope.SetColumn("CreatedOn", time.Now().Unix())

//     return nil
// }

// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//     scope.SetColumn("ModifiedOn", time.Now().Unix())

//     return nil
// }
