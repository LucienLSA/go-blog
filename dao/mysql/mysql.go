package mysql

import (
	"fmt"

	"github.com/LucienLSA/go-blog/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// 不要忘了导入数据库驱动

var db *sqlx.DB

func InitMysql(cfg *settings.MySQLConfig) (err error) {
	// dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		// viper.GetString("mysql.user"),
		// viper.GetString("mysql.password"),
		// viper.GetString("mysql.host"),
		// viper.GetInt("mysql.port"),
		// viper.GetString("mysql.dbname"),
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	// db = sqlx.MustConnect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	// db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	// db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	return
}

func Close() {
	_ = db.Close()
}
