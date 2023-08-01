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
	"xmlt/middle"
)

type ArticleRouter struct{ service service.ArticleService }

func (a *ArticleRouter) InitApiRouter(router *gin.RouterGroup) {
	articleRouter := router.Group("article")
	publicRouter := router.Group("article")
	articleRouter.Use(middle.VerifyJWTMiddleware())
	// 依赖注入
	artCache := cache.NewArticleRedisCache(global.Redis)
	a.service = service.NewArticleService(
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_ONLINE), artCache), // 加载线上库 - 用户
		repository.NewArticleRepo(dao.NewArticleDAO(global.DB_MAKE), artCache),   // 加载制作库 - 作者
	)
	articleCtl := v1.NewArticleController(a.service)
	{
		publicRouter.GET("/read/:id", articleCtl.GetByID)
		publicRouter.GET("/read/", articleCtl.GetArticleByCate)
		articleRouter.GET("/new", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "前端路由：新建文章")
		})
		articleRouter.GET("/update", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "前端路由：修改文章")
		})
		articleRouter.POST("/new", articleCtl.Save)
		articleRouter.POST("/publish", articleCtl.Publish)

	}
}
