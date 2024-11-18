package model

import "time"

type Model struct {
	CreatedOn  *time.Time `json:"created_on" gorm:"default:0;column:created_on" comment:"创建时间"`                // 创建时间
	CreatedBy  string     `json:"created_by" gorm:"default:'';column:created_by; varchar(100)" comment:"创建人"`  // 创建人
	ModifiedOn *time.Time `json:"modified_on" gorm:"default:0; column:modified_on" comment:"修改时间"`             // 修改时间             // 修改时间
	ModifiedBy string     `json:"modified_by" gorm:"default:'';varchar(100);column:modified_by" comment:"修改人"` // 修改人
	DeletedOn  *time.Time `json:"deleted_on" gorm:"default:0;column:deleted_on" comment:"删除时间"`                // 删除时间
	IsDel      uint8      `json:"is_del" gorm:"tinyint(3);default:0;column:is_del" comment:"是否删除 0-未删除 1-已删除"` // 是否删除 0-未删除 1-已删除
}
