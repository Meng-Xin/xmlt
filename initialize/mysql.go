package initialize

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"xmlt/global"
	"xmlt/internal/model"
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
	// 制作库开启 慢查询日志 Callback
	SlowQueryLog(dbMake)

	dbOnline, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       onlineDSN, // DSN data source name
		DefaultStringSize:         256,       // string 类型字段的默认长度
		DisableDatetimePrecision:  true,      // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,      // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,      // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,     // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	onlineDB, _ := dbOnline.DB()
	onlineDB.SetMaxIdleConns(20)  //设置连接池，空闲
	onlineDB.SetMaxOpenConns(100) //打开
	onlineDB.SetConnMaxLifetime(time.Second * 30)
	if err != nil {
		panic(err)
	}
	global.DB_MAKE = dbMake
	global.DB_ONLINE = dbOnline
	if global.Config.MysqlMake.AutoMigrate {
		Migration(global.DB_MAKE)
	}
	if global.Config.MysqlOnline.AutoMigrate {
		Migration(global.DB_ONLINE)
	}
}

// Migration 执行数据迁移
func Migration(db *gorm.DB) {
	//自动迁移模
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		model.User{},
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
				log.Debug("慢查询 %s", d.Statement.SQL.String())
			}
		}
	})
	if err != nil {
		panic(err)
	}
}
