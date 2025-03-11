package mysql

import (
	"context"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/settings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// 不要忘了导入数据库驱动

var (
	_db *gorm.DB
)

func InitMysql(cfg *settings.MySQLConfig) (err error) {
	pathRead := strings.Join([]string{
		cfg.User, ":", cfg.Password, "@tcp(", cfg.Host, ":", cfg.Port, ")/", cfg.DbName, "?charset=utf8mb4&parseTime=True"}, "")
	pathWrite := strings.Join([]string{
		cfg.User, ":", cfg.Password, "@tcp(", cfg.Host, ":", cfg.Port, ")/", cfg.DbName, "?charset=utf8mb4&parseTime=True"}, "")
	var ormLogger logger.Interface
	if gin.Mode() == settings.Conf.LogConfig.Level {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	_db = db
	_ = _db.Use(dbresolver.
		Register(dbresolver.Config{
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(pathRead)},                         // 读操作
			Replicas: []gorm.Dialector{mysql.Open(pathWrite), mysql.Open(pathWrite)}, // 写操作
			Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
		}))
	_db = _db.Set("gorm:table_options", "charset=utf8mb4")
	err = migrate()
	if err != nil {
		zap.L().Error("migrate failed, err:%v\n", zap.Error(err))
		return
	}
	return
}

// 创建数据库客户端上下文
// 用于将Go的context.Context对象与数据库操作绑定在一起，主要作用是对控制请求的超时、取消、追踪等。
func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
