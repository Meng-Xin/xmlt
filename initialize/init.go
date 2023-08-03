package initialize

import (
	"github.com/go-co-op/gocron"
	"time"
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
	// RabbitMQ初始化
	global.RabbitMQ = utils.NewRabbitMQ("comment", "", "", "amqp://guest:guest@localhost:5672/")
	// Cron 定时组件初始化
	global.Cron = gocron.NewScheduler(time.UTC)

}
