package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"xmlt/config"
)

var (
	Config    *config.AllConfig // 全局config
	DB_MAKE   *gorm.DB          // 制作库 - 作家|审核
	DB_ONLINE *gorm.DB          // 线上库 - 用户
	Redis     *redis.Client     // RedisClient
)
