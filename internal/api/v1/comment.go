package v1

import (
	"github.com/gin-gonic/gin"
	"xmlt/internal/service"
)

type CommentController struct {
	service service.CommentService
}

func NewCommentCtrl(commentService service.CommentService) *CommentController {
	return &CommentController{service: commentService}
}

func (c *CommentController) Create(ctx *gin.Context) {

}

func (c *CommentController) Delete(ctx *gin.Context) {

}
func (c *CommentController) Get(ctx *gin.Context) {

}
func (c *CommentController) GetByUserID(ctx *gin.Context) {

}
