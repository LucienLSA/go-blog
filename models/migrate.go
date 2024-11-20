package models

import (
	"blog-service/global"
)

func migrate() {
	err := global.DBEngine.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&Tag{},
			&Article{},
			&Author{},
		)
	if err != nil {
		panic(err)
	}
}
