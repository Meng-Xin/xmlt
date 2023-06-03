package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"xmlt/internal/domain"
	"xmlt/internal/service"
)

type ArticleController struct {
	service service.ArticleService
}

func NewArticleController(service service.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (h *ArticleController) GetByID(ctx *gin.Context) {

}

func (a *ArticleController) Save(ctx *gin.Context) {
	var vo ArticleVO
	if err := ctx.Bind(&vo); err != nil {
		// 出现 error 的情况下，实际上前端已经返回了
		return
	}
	// 缺乏登录部分，所以直接写死
	var authorID uint64 = 123
	_, err := a.service.Save(ctx.Request.Context(), domain.Article{
		Title:   vo.Title,
		Content: vo.Content,
		Author: domain.Author{
			ID: authorID,
		},
	})
	if err != nil {
		// 这边不能把 error 写回去
		// 暂时我直接输出到控制台
		log.Info(err)
		ctx.String(http.StatusInternalServerError, "系统异常，请重试")
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/article/new/success")
}

type ArticleVO struct {
	ID      uint64 `form:"id"`
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	// 一般来说，考虑到各种 APP 发版本不容易，
	// 所以数字、货币、日期、国际化之类的都是后端做的
	// 前端就是无脑展示
	Ctime string
	Utime string
}
