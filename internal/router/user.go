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

type UserRouter struct {
	service service.UserService
}

func (u *UserRouter) InitApiRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	// 依赖注入
	userCache := cache.NewUserRedisCache(global.Redis)
	userDao := dao.NewUserDao(global.DB_MAKE)
	u.service = service.NewUserService(
		repository.NewUserRepo(userDao, userCache),
	)
	userCtrl := v1.NewUserController(u.service)
	{
		userRouter.POST("/login", userCtrl.Login)
		userRouter.POST("/register", userCtrl.Register)
	}
}
