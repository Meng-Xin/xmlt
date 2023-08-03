package router

import (
	"github.com/gin-gonic/gin"
	v1 "xmlt/internal/api/v1"
)

type HomeRouter struct {
}

func (u *HomeRouter) InitApiRouter(router *gin.RouterGroup, all *RouteGroup) {
	homeRouter := router.Group("home")
	// 依赖注入
	userCtrl := v1.NewHomeController(
		all.UserRouter.service,
		all.ArticleRouter.service,
		all.CategoryRouter.service,
		all.UserLikeArticleRouter.service)
	{
		homeRouter.GET("", userCtrl.GetHomeInfo)
	}
}
