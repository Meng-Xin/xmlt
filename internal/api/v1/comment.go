package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
)

type CommentController struct {
	service service.CommentService
}

func NewCommentCtrl(commentService service.CommentService) *CommentController {
	return &CommentController{service: commentService}
}

func (c *CommentController) Create(ctx *gin.Context) {
	code := e.SUCCESS
	var com comment
	err := ctx.Bind(&com)
	if err != nil {
		global.Log.Info(err.Error())
		return
	}
	uid := ctx.GetUint64(enum.CtxUid)
	// 新增评论
	err = c.service.AddComment(ctx, domain.Comment{
		Content:   com.Content,
		UserID:    uid,
		ArticleID: com.ArticleID,
		ParentID:  com.ParentID,
	})
	if err != nil {
		global.Log.Info(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Msg:    "新增评论成功",
	})
}

func (c *CommentController) Delete(ctx *gin.Context) {

}
func (c *CommentController) Get(ctx *gin.Context) {

}
func (c *CommentController) GetByUserID(ctx *gin.Context) {

}

type comment struct {
	ID        uint64 `json:"id" form:"id"`           // 评论ID
	Content   string `json:"content" form:"content"` // 评论内容
	UserID    uint64 `json:"userID"`                 // 评论用户ID
	ArticleID uint64 `json:"articleID"`              // 文章ID
	ParentID  uint64 `json:"parent_id"`              // 父级评论ID
	Floor     uint32 `json:"floor"`                  // 评论楼层
	State     uint8  `json:"state"`                  // 该评论状态 0:正常，1：删除

	Ctime int64 `json:"ctime"` // 创建时间，毫秒作为单位
	Utime int64 `json:"utime"` // 更新时间，毫秒作为单位
}
