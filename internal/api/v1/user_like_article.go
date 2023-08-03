package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xmlt/global"
	"xmlt/internal/service"
	"xmlt/internal/shared/e"
	"xmlt/internal/shared/enum"
	"xmlt/internal/shared/public"
)

type UserLikeArticleController struct {
	service service.UserLikeArticleService
}

func NewUserLikeArticleController(service service.UserLikeArticleService) *UserLikeArticleController {
	return &UserLikeArticleController{service: service}
}

func (u *UserLikeArticleController) LikeOrCancelLike(ctx *gin.Context) {
	code := e.SUCCESS
	var vo vo
	err := ctx.Bind(&vo)
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Info(err.Error())
		return
	}
	uid := ctx.GetUint64(enum.CtxUid)
	err = u.service.UpdateLike(ctx, vo.ArticleID, uint64(uid))
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Info(err.Error())
		return
	}
	ctx.JSON(code, public.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	})
}
func (u *UserLikeArticleController) GetLikeState(ctx *gin.Context) {
	code := e.SUCCESS
	var data vo
	err := ctx.Bind(&data)
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Info(err.Error())
		return
	}
	uid := ctx.GetUint64(enum.CtxUid)
	state := u.service.GetLikeState(ctx, data.ArticleID, uid)
	data.UserID = uid
	data.LikeState = state
	ctx.JSON(code, public.Response{
		Status: code,
		Data:   data,
	})
}
func (u *UserLikeArticleController) GetArticleLikeNum(ctx *gin.Context) {
	code := e.SUCCESS

	articleID := ctx.Param("id")
	parseUint, err := strconv.ParseUint(articleID, 10, 64)
	if err != nil {
		return
	}
	var data ArticleLike
	likeNum, err := u.service.GetLikeNum(ctx, parseUint)
	data.ArticleID = parseUint
	data.LikeNum = likeNum
	ctx.JSON(code, public.Response{
		Status: code,
		Data:   data,
	})
}

type vo struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	ArticleID uint64 `json:"article_id"`
	LikeState bool   `json:"like_state"`

	Ctime uint64 `json:"ctime"`
	Utime uint64 `json:"utime"`
}

type ArticleLike struct {
	ArticleID uint64 `json:"article_id"`
	LikeNum   uint64 `json:"like_num"`
	LikeUser  []uint64
}
