package global

import (
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"xmlt/config"
	"xmlt/utils"
)

var (
	Config   *config.AllConfig // 全局config
	DB_MAKE  *gorm.DB          // 制作库 - 作家|审核
	Redis    *redis.Client     // RedisClient
	Log      utils.ILog        // 异常日志系统
	RabbitMQ *utils.RabbitMQ
	Cron     *gocron.Scheduler
)
