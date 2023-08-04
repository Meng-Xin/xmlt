package main

import (
	"net/http"
	_ "runtime/pprof"
	"xmlt/global"
	"xmlt/initialize"
	"xmlt/internal/shared/crontask"
)

func main() {
	// 加载初始化文件
	initialize.GlobalInit()
	// 开启pprof
	go http.ListenAndServe("9191", nil)

	// 注册路由
	r := initialize.NewRouter()

	// 开启定时任务
	cron := crontask.NewUserLikeCron(global.Cron, global.Redis, global.DB_MAKE)
	cron.SyncFromRedisToMysql()
	global.Cron.StartAsync()

	r.Run(global.Config.Server.DSN())
}
