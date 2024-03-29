package initialize

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"xmlt/global"
	"xmlt/internal/model"
)

var (
	GormToManyRequestError = errors.New("gorm: to many request")
)

func InitDatabase(makeDSN string, onlineDSN string) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	dbMake, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       makeDSN, // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	sqlDB, _ := dbMake.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	// 中间件引入
	SlowQueryLog(dbMake)
	GormRateLimiter(dbMake, rate.NewLimiter(500, 1000))
	if err != nil {
		panic(err)
	}
	global.DB_MAKE = dbMake
	if global.Config.MysqlMake.AutoMigrate {
		Migration(global.DB_MAKE)
	}
}

// Migration 执行数据迁移
func Migration(db *gorm.DB) {
	//自动迁移模
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		model.User{},
		model.Article{},
		model.UserLikeArticle{},
		model.Tag{},
		model.Category{},
		model.Comment{},
		model.Notification{},
	)
	log.Info("register table success")
}

// SlowQueryLog 慢查询日志
func SlowQueryLog(db *gorm.DB) {
	err := db.Callback().Query().Before("*").Register("slow_query_start", func(d *gorm.DB) {
		now := time.Now()
		d.Set("start_time", now)
	})
	if err != nil {
		panic(err)
	}

	err = db.Callback().Query().After("*").Register("slow_query_end", func(d *gorm.DB) {
		now := time.Now()
		start, ok := d.Get("start_time")
		if ok {
			duration := now.Sub(start.(time.Time))
			// 一般认为 10 Ms 为Sql慢查询
			if duration > time.Millisecond*10 {
				global.Log.Error("慢查询 %s", d.Statement.SQL.String())
			}
		}
	})
	if err != nil {
		panic(err)
	}
}

// GormRateLimiter Gorm限流器 此限流器不能终止GORM查询链。
func GormRateLimiter(db *gorm.DB, r *rate.Limiter) {
	err := db.Callback().Query().Before("*").Register("RateLimitGormMiddleware", func(d *gorm.DB) {
		if !r.Allow() {
			d.AddError(GormToManyRequestError)
			global.Log.Error(GormToManyRequestError.Error())
			return
		}
	})
	if err != nil {
		panic(err)
	}
}
