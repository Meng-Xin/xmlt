package main

import (
	"net/http"
	"xmlt/global"
	"xmlt/initialize"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/router"
	"xmlt/internal/service"
)

func main() {
	// 加载初始化文件
	initialize.GlobalInit()
	// 开启pprof
	go http.ListenAndServe("9191", nil)
	// 依赖注入
	artCache := cache.NewArticleRedisCache(global.Redis)
	articleService := service.NewArticleService(
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_ONLINE), nil),    // 加载线上库 - 用户
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_MAKE), artCache), // 加载制作库 - 作者
	)
	v1.NewArticleController(articleService)
	// 注册路由
	r := router.NewRouter()
	r.Run(global.Config.Server.DSN())
}
