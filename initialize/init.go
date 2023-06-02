package initialize

import (
	"xmlt/config"
	"xmlt/global"
)

func GlobalInit() {
	// 初始化全局配置文件
	global.Config = config.InitLoadConfig()
	// gorm 初始化
	InitDatabase(global.Config.MysqlMake.Dsn(), global.Config.MysqlOnline.Dsn())
	// cache 初始化
	InitRedis()
}
