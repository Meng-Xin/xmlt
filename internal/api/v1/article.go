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

type ArticleController struct {
	service service.ArticleService
}

func NewArticleController(service service.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (a *ArticleController) GetByID(ctx *gin.Context) {
	code := e.SUCCESS
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		code = e.ArticleDoseNotExist
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	art, err := a.service.Get(ctx.Request.Context(), id, enum.ArticleGetSourceOnline)
	if err != nil {
		// TODO 如果代码严谨的话，这边要区别是真的没有数据，还是服务器出现了异常
		return
	}
	// 构建相应信息
	var articleRes ArticleRes
	articleRes.init(art)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   articleRes,
	})
}

func (a *ArticleController) GetArticleByCate(ctx *gin.Context) {
	code := e.SUCCESS
	category, _ := strconv.ParseUint(ctx.Query("category"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	pageEntity := shared.NewPage(page, pageSize)
	arts, err := a.service.GetCategoryArticles(ctx.Request.Context(), category, pageEntity)
	if err != nil {
		// TODO 如果代码严谨的话，这边要区别是真的没有数据，还是服务器出现了异常
		return
	}

	// 构建相应信息
	var articlesRes ArticlesRes
	articlesRes.init(arts, pageEntity)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   articlesRes,
	})
}

// Save 作者可以保存文章
func (a *ArticleController) Save(ctx *gin.Context) {
	code := e.SUCCESS
	var data ArticleSaveReq
	if err := ctx.Bind(&data); err != nil {
		// 出现 error 的情况下，实际上前端已经返回了
		global.Log.Warn("解析Article Save 结构体错误：", err.Error())
		return
	}
	// 从 jwtmiddler 中拿到解析的UID
	authorID := ctx.GetUint64(enum.CtxUid)
	artId, err := a.service.Save(ctx.Request.Context(), domain.Article{
		Title:      data.Title,
		Content:    data.Content,
		CategoryID: data.CategoryID,
		Author:     authorID,
	})
	if err != nil {
		// 这边不能把 error 写回去
		// 暂时我直接输出到控制台
		global.Log.Error(err)
		ctx.String(http.StatusInternalServerError, "系统异常，请重试")
		return
	}

	var articleSaveRes ArticleSaveRes
	articleSaveRes.init(artId)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   articleSaveRes,
		Msg:    "保存成功！",
	})
}

func (a *ArticleController) Publish(ctx *gin.Context) {
	// TODO 跨领域调用，作者发布后需要通过,Notification模块进行消息推送。
	code := e.SUCCESS
	var data ArticlePublishReq
	err := ctx.Bind(&data)
	if err != nil {
		code = e.ERROR
		global.Log.Warn("解析Article Save 结构体错误：", err.Error())
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	authorID := ctx.GetUint64(enum.CtxUid)
	publishId, err := a.service.Publish(ctx, data.ArticleId, authorID)
	if err != nil {
		code = e.ERROR
		global.Log.Warn(err.Error())
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	var articlePublishRes ArticlePublishRes
	articlePublishRes.init(publishId)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   articlePublishRes,
		Msg:    "发布成功，审核中……",
	})
}

// ArticleSaveReq POST:/article/new
type ArticleSaveReq struct {
	Title      string `form:"title" json:"title" binding:"required"`           // 帖子标题
	Content    string `form:"content" json:"content" binding:"required"`       // 帖子内容
	CategoryID uint64 `form:"categoryID" json:"categoryID" binding:"required"` // 所属板块
}

type ArticleSaveRes struct {
	ArticleId uint64 `from:"article_id" json:"id"`
}

func (a *ArticleSaveRes) init(artId uint64) {
	a.ArticleId = artId
}

type ArticlePublishReq struct {
	ArticleId uint64 `form:"article_id" json:"article_id"`
}

type ArticlePublishRes struct {
	ArticleId uint64
}

func (a *ArticlePublishRes) init(artId uint64) {
	a.ArticleId = artId
}

// ArticleRes GetByID -> Get:article/:id
type ArticleRes struct {
	ID          uint64   // 帖子ID
	Title       string   // 帖子标题
	Content     string   // 帖子内容
	Status      uint8    // 帖子状态
	CategoryID  uint64   // 所属板块
	NiceTopic   uint8    // 精选话题
	BrowseCount uint64   // 浏览量
	ThumbsUP    uint64   // 点赞数
	Author      struct { // 作者信息
		ID       uint64
		NickName string
		Avatar   string
	}
	Ctime string
	Utime string
}

func (a *ArticleRes) init(art domain.Article) {
	a.ID = art.ID
	a.Title = art.Title
	a.Content = art.Content
	a.Status = art.Status
	a.CategoryID = art.CategoryID
	a.NiceTopic = art.NiceTopic
	a.BrowseCount = art.BrowseCount
	a.ThumbsUP = art.ThumbsUP
	a.Author.ID = art.User.ID
	a.Author.NickName = art.User.NickName
	a.Author.Avatar = art.User.Avatar
	a.Ctime = art.Ctime.String()
	a.Utime = art.Utime.String()
}

// ArticlesRes GetArticleByCate -> Get：article/read?category=
type ArticlesRes struct {
	PageInfo *shared.Page // 分页信息
	Articles []ArticleRes // 文章列表
}

func (a *ArticlesRes) init(arts []domain.Article, page *shared.Page) {
	for _, art := range arts {
		var article ArticleRes
		article.init(art)
		a.Articles = append(a.Articles, article)
	}
	a.PageInfo = page
}
