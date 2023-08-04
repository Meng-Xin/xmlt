package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"xmlt/global"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/service"
	"xmlt/middle"
)

type CommentRouter struct {
	service service.CommentService
}

func (c *CommentRouter) InitApiRouter(router *gin.RouterGroup) {
	commentRouter := router.Group("comment")
	pubCommentRouter := router.Group("comment")
	// 使用中间件
	commentRouter.Use(middle.VerifyJWTMiddleware())

	// 依赖注入
	commentCache := cache.NewCommentRedisCache(global.Redis)
	commentDao := dao.NewCommentDao(global.DB_MAKE)
	commentRepo := repository.NewCommentRepo(commentDao, commentCache)
	c.service = service.NewCommentService(commentRepo)

	// 挂载 RabbitMQ 对创建评论的处理
	err := commentRepo.ConsumerMQ(context.Background())

	if err != nil {
		global.Log.Error(err.Error())
	}
	commentCtrl := v1.NewCommentCtrl(c.service)
	{
		// 新增评论
		commentRouter.POST("/new", commentCtrl.Create)
		// 删除评论
		commentRouter.POST("/delete", commentCtrl.Delete)
		// 拉取文章下的所有评论
		pubCommentRouter.GET("/get", commentCtrl.GetArtComment)
		// 拉取某个用户所有评论
		commentRouter.POST("/getByUserID", commentCtrl.GetByUserID)
		// TODO 拉取某个单独评论信息，这个一般是在通知模块用的，这里写这个接口不合适
		//commentRouter.POST("/getByCommentID")
	}
}
