package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmlt/global"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/service"
)

type ArticleRouter struct{}

func (a *ArticleRouter) InitApiRouter(router *gin.RouterGroup) {
	articleRouter := router.Group("article")
	// 依赖注入
	artCache := cache.NewArticleRedisCache(global.Redis)
	articleService := service.NewArticleService(
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_ONLINE), nil),    // 加载线上库 - 用户
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_MAKE), artCache), // 加载制作库 - 作者
	)
	articleCtl := v1.NewArticleController(articleService)
	{
		articleRouter.GET("/read/:id", articleCtl.GetByID)
		articleRouter.GET("/new", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "新建文章")
		})
		articleRouter.POST("/new", articleCtl.Save)
		articleRouter.POST("/publish")
	}
}
