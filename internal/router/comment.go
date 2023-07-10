package router

import (
	"github.com/gin-gonic/gin"
	"xmlt/global"
	v1 "xmlt/internal/api/v1"
	"xmlt/internal/repository"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/service"
	"xmlt/middle"
)

type CommentRouter struct{}

func (a *CommentRouter) InitApiRouter(router *gin.RouterGroup) {
	commentRouter := router.Group("comment")
	// 使用中间件
	commentRouter.Use(middle.VerifyJWTMiddleware())
	// 依赖注入
	commentCache := cache.NewCommentRedisCache(global.Redis)
	commentDao := dao.NewCommentDao(global.DB_MAKE)
	commentService := service.NewCommentService(repository.NewCommentRepo(commentDao, commentCache))

	commentCtrl := v1.NewCommentCtrl(commentService)
	{
		// 新增评论
		commentRouter.POST("/new", commentCtrl.Create)
		// 删除评论
		commentRouter.POST("/delete", commentCtrl.Delete)
		// 拉取文章下的所有评论
		commentRouter.POST("/get", commentCtrl.Get)
		// 拉取某个用户所有评论
		commentRouter.POST("/getByUserID", commentCtrl.GetByUserID)
		// TODO 拉取某个单独评论信息，这个一般是在通知模块用的，这里写这个接口不合适
		//commentRouter.POST("/getByCommentID")
	}
}
