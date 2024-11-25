package models

import (
	"blog-service/conf"
	"blog-service/global"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySQLEngine(mysqlSetting *conf.MysqlSettingS) (*gorm.DB, error) {
	// conn := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true"
	// dsn := fmt.Sprintf(conn,
	// 	databaseSetting.UserName,
	// 	databaseSetting.Password,
	// 	databaseSetting.Host,
	// 	databaseSetting.Port,
	// 	databaseSetting.DBName,
	// 	databaseSetting.Charset,
	// )

	dsn := strings.Join([]string{mysqlSetting.UserName, ":", mysqlSetting.Password, "@tcp(", mysqlSetting.Host, ":", mysqlSetting.Port,
		")/", mysqlSetting.DBName, "?charset=", mysqlSetting.Charset, "&parseTime=true"}, "")
	fmt.Println(dsn)
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
			TablePrefix:   "blog_", //表名前缀，`Tag` 的表名应该是 `blog_tag`
			SingularTable: true,    // 表名不加s
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlSetting.MaxOpenConns)
	global.DBEngine = db
	migrate()
	return db, nil
}
