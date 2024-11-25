package conf

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	vp *viper.Viper
}

func InitConfig() (*Config, error) {
	// 获取当前工作目录的路径
	workDir, _ := os.Getwd()
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	// 将配置路径添加为当前工作目录下的 config/local 目录
	vp.AddConfigPath(workDir + "/conf/local")
	// 将当前工作目录也添加为可能的配置文件路径
	vp.AddConfigPath(workDir)
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Config{vp}, nil
}
