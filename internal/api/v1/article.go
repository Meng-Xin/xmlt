package v1

import (
	"github.com/gin-gonic/gin"
	"xmlt/internal/service"
)

type ArticleController struct {
	service service.ArticleService
}

func NewArticleController(service service.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (h *ArticleController) Get(ctx *gin.Context) {

}
