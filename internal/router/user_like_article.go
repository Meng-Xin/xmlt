package router

import (
	"github.com/gin-gonic/gin"
	"xmlt/global"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/service"
	"xmlt/middle"
)

type UserLikeArticleRouter struct{}

func (a *UserLikeArticleRouter) InitApiRouter(router *gin.RouterGroup) {
	userLikeArtRouter := router.Group("userLikeArticle")
	userLikeArtRouter.Use(middle.VerifyJWTMiddleware())
	// 依赖注入
	userLikeArtCache := cache.NewUserLikeArticleCache(global.Redis)
	userService := service.NewUserLikeArticleService(
		repository.NewUserLikeArticleRepo(userLikeArtCache),
	)
	userLikeCtrl := v1.NewUserLikeArticleController(userService)
	{
		userLikeArtRouter.POST("/likeOrCancelLike", userLikeCtrl.LikeOrCancelLike)
		userLikeArtRouter.POST("/getUserLikeState", userLikeCtrl.GetLikeState)
		userLikeArtRouter.GET("/getArticleLikeNum/:id", userLikeCtrl.GetArticleLikeNum)
	}
}
