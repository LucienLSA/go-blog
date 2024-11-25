package global

import "blog-service/conf"

// 定义全局变量 配置信息
var (
	ServerSetting *conf.ServerSettingS
	AppSetting    *conf.AppSettingS
	MysqlSetting  *conf.MysqlSettingS
	RedisSetting  *conf.RedisSettingS
)
