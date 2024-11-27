package conf

import (
	"blog-service/pkg/logging"
	"fmt"
	"log"
	"reflect"
	"time"
)

type ServerSettingS struct {
	RunMode      string        `yaml:"runMode"`
	HttpPort     string        `yaml:"httpPort"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`  // 读取超时时间
	WriteTimeout time.Duration `yaml:"writeTimeout"` // 写入超时时间
}

type AppSettingS struct {
	RuntimeRootPath string `yaml:"runtimeRootPath"`
	DefaultPageSize int    `yaml:"defaultPageSize"` // 默认每页显示的记录数
	MaxPageSize     int    `yaml:"maxPageSize"`     // 最大每页显示的记录数
	LogSavePath     string `yaml:"logSavePath"`
	LogFileName     string `yaml:"logFileName"` // 日志文件名
	LogFileExt      string `yaml:"logFileExt"`  // 日志文件扩展名

	ImgSavePath    string   `yaml:"imgSavePath"`
	ImageMaxSize   int      `yaml:"imageMaxSize"`   // 图片最大大小（单位：MB）
	ImageAllowExts []string `yaml:"imageAllowExts"` // 允许的图片扩展名

	JwtSecret string `yaml:"jwtSecret"` // JWT 密钥

	PrefixUrl      string `yaml:"prefixUrl"`      // URL前缀
	ExportSavePath string `yaml:"exportSavePath"` // 导出文件保存路径
}

type MysqlSettingS struct {
	DBType       string `yaml:"dBType"`       // 数据库类型
	UserName     string `yaml:"userName"`     // 数据库用户名
	Password     string `yaml:"password"`     // 数据库密码
	Host         string `yaml:"host"`         // 数据库主机地址
	Port         string `yaml:"port"`         // 数据库端口
	DBName       string `yaml:"dBName"`       // 数据库名称
	TablePrefix  string `yaml:"tablePrefix"`  // 数据库表前缀
	Charset      string `yaml:"charset"`      // 数据库字符集
	MaxIdleConns int    `yaml:"maxIdleConns"` // 最大空闲连接数
	MaxOpenConns int    `yaml:"maxOpenConns"` // 最大打开连接数
}

type RedisSettingS struct {
	DBName      int           `yaml:"dBName"`      // Redis 数据库编号
	Host        string        `yaml:"host"`        // Redis 主机地址
	Port        string        `yaml:"port"`        // Redis 端口
	Password    string        `yaml:"password"`    // Redis 密码
	MaxIdle     int           `yaml:"maxIdle"`     // 最大空闲连接数
	MinIdle     int           `yaml:"minIdle"`     // 最小空闲连接数
	MaxActive   int           `yaml:"maxActive"`   // 最大连接数
	IdleTimeout time.Duration `yaml:"idleTimeout"` // 空闲连接超时时间（秒）
}

// type ServerSettingS struct {
// 	RunMode      string
// 	HttpPort     string
// 	ReadTimeout  time.Duration
// 	WriteTimeout time.Duration
// }

// type AppSettingS struct {
// 	RuntimeRootPath string
// 	DefaultPageSize int
// 	MaxPageSize     int
// 	LogSavePath     string
// 	LogFileName     string
// 	LogFileExt      string

// 	ImgSavePath    string
// 	ImageMaxSize   int
// 	ImageAllowExts []string
// 	ImagePrefixUrl string
// 	JwtSecret      string
// }

// type DatabaseSettingS struct {
// 	DBType       string
// 	UserName     string
// 	Password     string
// 	Host         string
// 	Port         string
// 	DBName       string
// 	TablePrefix  string
// 	Charset      string
// 	MaxIdleConns int
// 	MaxOpenConns int
// }

func (c *Config) ReadSection(k string, v interface{}) error {
	err := c.vp.UnmarshalKey(k, v)
	if err != nil {
		log.Fatalf("unable to decode into config struct, %v", err)
		logging.LogrusObj.Infoln("配置文件读取失败", err)
		return err
	}
	fmt.Println("配置文件读取成功")
	// fmt.Println("配置文件内容:", v)
	value := reflect.ValueOf(v).Elem()
	fmt.Println("反射值:", value)
	fmt.Println("反射类型:", value.Type())
	return nil
}
