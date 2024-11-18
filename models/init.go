package model

import (
	"blog-service/global"
	"blog-service/pkg/setting"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	conn := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(conn,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.Port,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)

	var ormLogger logger.Interface
	if global.ServerSetting.RunMode == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加s
		},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	global.DBEngine = db
	migrate()
	return db, nil
}

// func NewDBEngine(ctx context.Context) *gorm.DB {
// 	db := global.DBEngine
// 	return db.WithContext(ctx)
// }
