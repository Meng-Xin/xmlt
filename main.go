package main

import (
	"net/http"
	_ "runtime/pprof"
	"xmlt/global"
	"xmlt/initialize"
)

func main() {
	// 加载初始化文件
	initialize.GlobalInit()
	// 开启pprof
	go http.ListenAndServe("9191", nil)

	// 注册路由
	r := initialize.NewRouter()
	r.Run(global.Config.Server.DSN())
}
