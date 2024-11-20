package setting

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	vp *viper.Viper
}

// 通过viper读取配置文件
func NewConfig() (*Config, error) {
	workDir, _ := os.Getwd()
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(workDir + "/conf/")
	vp.AddConfigPath(workDir)
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Config{vp}, nil
}
