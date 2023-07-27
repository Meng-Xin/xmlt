package router

import (
	"github.com/gin-gonic/gin"
	"xmlt/global"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/service"
)

type CategoryRouter struct{}

func (a *CategoryRouter) InitApiRouter(router *gin.RouterGroup) {
	categoryRouter := router.Group("category")
	// 依赖注入
	categoryDao := dao.NewCategoryDao(global.DB_MAKE)
	categoryCache := cache.NewCategoryCache(global.Redis)
	categoryService := service.NewCategoryService(
		repository.NewCategoryRepo(categoryDao, categoryCache),
	)
	categoryCtrl := v1.NewCategoryController(categoryService)
	{
		categoryRouter.POST("/create", categoryCtrl.CreateCategory)
		categoryRouter.DELETE("/delete", categoryCtrl.DeleteCategory)
		categoryRouter.PUT("/update", categoryCtrl.UpdateCategory)
		categoryRouter.GET("/get", categoryCtrl.GetCategoryList)
	}
}
