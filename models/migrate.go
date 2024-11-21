package models

import (
	"blog-service/global"
	"blog-service/pkg/logging"
)

func migrate() {
	err := global.DBEngine.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&Tag{},
			&Article{},
			&Author{},
		)
	if err != nil {
		logging.LogrusObj.Panicln(err)
		panic(err)
	}
}
