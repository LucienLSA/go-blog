package models

import "time"

type Model struct {
	ID         int64      `json:"id" gorm:"column:id; primary_key;not null;auto_increment"`
	CreatedOn  *time.Time `json:"created_on" gorm:"column:created_on"`   // 创建时间
	ModifiedOn *time.Time `json:"modified_on" gorm:"column:modified_on"` // 修改时间             // 修改时间
	DeletedOn  *time.Time `json:"deleted_on" gorm:"column:deleted_on"`   // 删除时间             // 删除时间
}
