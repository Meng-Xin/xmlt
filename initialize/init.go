package initialize

import (
	"xmlt/config"
	"xmlt/global"
	"xmlt/utils"
)

func GlobalInit() {
	// 初始化全局配置文件
	global.Config = config.InitLoadConfig()
	// 日志系统初始化
	global.Log = utils.NewLLogger()
	// gorm 初始化
	InitDatabase(global.Config.MysqlMake.Dsn(), global.Config.MysqlOnline.Dsn())
	// cache 初始化
	InitRedis()
}
