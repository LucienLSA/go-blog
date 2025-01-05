package settings

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         string `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	FilePath   string `mapstructure:"filepath"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// func InitSettings() (err error) {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	// viper.SetConfigFile("./config/config.yaml")
// 	viper.AddConfigPath("./config/")
// 	viper.AddConfigPath("./")
// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		// 读取配置信息失败
// 		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
// 		return
// 	}
// 	// 配置信息的反序列化
// 	if err := viper.Unmarshal(Conf); err != nil {
// 		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
// 	}
// 	viper.WatchConfig()
// 	viper.OnConfigChange(func(in fsnotify.Event) {
// 		fmt.Println("配置文件被修改了")
// 		if err := viper.Unmarshal(Conf); err != nil {
// 			fmt.Printf("viper Unmarshal failed, err:%v\n", err)
// 		}
// 	})
// 	return
// }

func InitSettings() (err error) {
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.SetConfigFile("./config/config.yaml")
	// viper.AddConfigPath("./config/")
	// viper.AddConfigPath("./")

	// viper.SetConfigFile(filePath)
	var filePath string
	flag.StringVar(&filePath, "filePath", "./config/config.yaml", "配置文件")
	//解析命令行参数
	flag.Parse()
	fmt.Println(filePath)
	fmt.Println(flag.Args())
	fmt.Println(flag.NArg())
	fmt.Println(flag.NFlag())
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}
	// 配置信息的反序列化
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
