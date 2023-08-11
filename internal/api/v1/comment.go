package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/service"
	"xmlt/internal/shared"
	"xmlt/internal/shared/e"
	"xmlt/internal/shared/enum"
	"xmlt/internal/shared/public"
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

func (c *CommentController) GetArtComment(ctx *gin.Context) {
	code := e.SUCCESS
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	articleId, _ := strconv.ParseUint(ctx.Query("article_id"), 10, 64)

	pageInfo := shared.NewPage(page, pageSize)
	rangeBy := shared.NewRange(pageInfo)
	artComment, err := c.service.GetArtComment(ctx.Request.Context(), articleId, pageInfo, *rangeBy)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   artComment,
	})
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

type CreateCommentReq struct {
	Content   string `json:"content" form:"content"`   // 评论内容
	ArticleID uint64 `json:"articleID" form:"comment"` // 文章ID
	ParentID  uint64 `json:"parent_id" form:"comment"` // 父级评论ID
}

type GetArtCommentRes struct {
	ID      uint64   // 评论ID
	Content string   // 评论内容
	User    struct { // 评论用户
		Uid      uint64
		NickName string
		Avatar   string
	}
	ArticleID struct { // 评论文章
		ArtId uint64
		Title string
	}
	ParentComment struct { // 父级评论信息
		ParentId uint64
		Content  string
	}
}
