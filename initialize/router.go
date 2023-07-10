package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xmlt/global"
	"xmlt/internal/api"
	"xmlt/internal/router"
	"xmlt/middle"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	allRouter := router.RouterGroupCtrl

	// 全局限流
	r.Use(middle.NewRateLimiterMiddleware(global.Redis, "general", 200, 60*time.Second))
	// 全局日志中间件
	r.Use(middle.LoggerToFile())

	var publicGroup, privateGroup, v1 *gin.RouterGroup
	// 公共路由组
	publicGroup = r.Group("/")
	{
		publicGroup.GET("/ping", api.Ping)
	}
	// 私有路由组
	privateGroup = r.Group("admin")
	privateGroup.Use(middle.VerifyJWTMiddleware())
	{
		privateGroup.GET("/login", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "JWT 校验通过！")
		})
	}
	// Swagger

	// V1管理
	v1 = r.Group("/api/v1")
	{
		// Article 帖子
		allRouter.ArticleRouter.InitApiRouter(v1)
		// User 用户
		allRouter.UserRouter.InitApiRouter(v1)
		// Comment 评论
		allRouter.CommentRouter.InitApiRouter(v1)
	}

	return r
}
