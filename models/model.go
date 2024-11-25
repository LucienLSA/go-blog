package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID         int64          `json:"id" gorm:"column:id; primary_key;not null;auto_increment"`
	CreatedOn  *time.Time     `json:"created_on" gorm:"column:created_on"`   // 创建时间
	ModifiedOn *time.Time     `json:"modified_on" gorm:"column:modified_on"` // 修改时间             // 修改时间
	DeletedOn  gorm.DeletedAt `json:"deleted_on" gorm:"column:deleted_on"`   // 删除时间             // 删除时间
}

// // updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
// func updateTimeStampForCreateCallback(scope *gorm.Scope) {
// 	if !scope.HasError() {
// 		nowTime := time.Now().Unix()
// 		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
// 			if createTimeField.IsBlank {
// 				createTimeField.Set(nowTime)
// 			}
// 		}

// 		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
// 			if modifyTimeField.IsBlank {
// 				modifyTimeField.Set(nowTime)
// 			}
// 		}
// 	}
// }

// BeforeCreate 设置创建和修改时间戳
func (model *Model) BeforeCreate(tx *gorm.DB) (err error) {
	nowTime := time.Now()
	if model.CreatedOn == nil { // 检查 CreatedOn 是否为空
		model.CreatedOn = &nowTime
	}
	return nil
}

// // updateTimeStampForUpdateCallback will set `ModifyTime` when updating
// func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
// 	if _, ok := scope.Get("gorm:update_column"); !ok {
// 		scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	}
// }

func (model *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	nowTime := time.Now()
	if model.ModifiedOn == nil { // 检查 ModifiedOn 是否为空
		model.ModifiedOn = &nowTime
	}
	return nil
}

// 软删除
