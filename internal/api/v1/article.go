package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
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
	var vo ArticleSave
	vo.init(art)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   vo,
	})
}

func (a *ArticleController) GetArticleByCate(ctx *gin.Context) {
	code := e.SUCCESS
	category, _ := strconv.ParseUint(ctx.Query("category"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	pageEntity := domain.NewPage(page, pageSize)
	arts, err := a.service.GetCategoryArticles(ctx.Request.Context(), category, pageEntity)
	if err != nil {
		// TODO 如果代码严谨的话，这边要区别是真的没有数据，还是服务器出现了异常
		return
	}

	var vo []ArticleSave
	for i, _ := range arts {
		articleVo := ArticleSave{}
		articleVo.init(arts[i])
		vo = append(vo, articleVo)
	}

	response := ArticleListResponse{
		PageInfo: pageEntity,
		Articles: vo,
	}
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   response,
	})
}

// Save 作者可以保存文章
func (a *ArticleController) Save(ctx *gin.Context) {
	code := e.SUCCESS
	var data ArticleSave
	if err := ctx.Bind(&data); err != nil {
		// 出现 error 的情况下，实际上前端已经返回了
		global.Log.Warn("解析Article Save 结构体错误：", err.Error())
		return
	}
	// 从 jwtmiddler 中拿到解析的UID
	authorID := ctx.GetUint64(enum.CtxUid)
	_, err := a.service.Save(ctx.Request.Context(), domain.Article{
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
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Msg:    "保存成功！",
	})
}

func (a *ArticleController) Publish(ctx *gin.Context) {
	code := e.SUCCESS

	var data ArticleSave
	err := ctx.Bind(&data)
	if err != nil {
		global.Log.Warn("解析Article Save 结构体错误：", err.Error())
		return
	}
	// CheckPushPermission 检查权限
	authorID := ctx.GetUint64(enum.CtxUid)
	err = a.service.Publish(ctx, domain.Article{
		ID:         data.ID,
		Title:      data.Title,
		Content:    data.Content,
		CategoryID: data.CategoryID,
		Author:     authorID,
	})
	if err != nil {
		if err == e.UpdateArticleIdenticalError {
			global.Log.Warn("非法修改文章作者身份,记录IP：", ctx.Request.Host)
		}
		return
	}
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Msg:    "发布成功，审核中……",
	})
}

type ArticleSave struct {
	ID          uint64 `form:"id" json:"id"`
	Title       string `form:"title" json:"title" binding:"required"`           // 帖子标题
	Content     string `form:"content" json:"content" binding:"required"`       // 帖子内容
	Status      uint8  `form:"status" json:"status"`                            // 帖子状态
	CategoryID  uint64 `json:"categoryID" form:"categoryID" binding:"required"` // 所属板块
	NiceTopic   uint8  `json:"niceTopic" from:"niceTopic"`                      // 精选话题
	BrowseCount uint64 `json:"browseCount" from:"browseCount"`                  // 浏览量
	ThumbsUP    uint64 `json:"thumbsUP" form:"thumbsUP"`                        // 点赞数
	Author      uint64 // 作者
	// 一般来说，考虑到各种 APP 发版本不容易，
	// 所以数字、货币、日期、国际化之类的都是后端做的
	// 前端就是无脑展示
	Ctime string
	Utime string
}

func (a *ArticleSave) init(art domain.Article) {
	a.ID = art.ID
	a.Title = art.Title
	a.Content = art.Content
	a.Status = art.Status
	a.CategoryID = art.CategoryID
	a.NiceTopic = art.NiceTopic
	a.BrowseCount = art.BrowseCount
	a.ThumbsUP = art.ThumbsUP
	a.Author = art.Author
	a.Ctime = art.Ctime.String()
	a.Utime = art.Utime.String()
}

type ArticleListResponse struct {
	PageInfo *domain.Page  // 分页信息
	Articles []ArticleSave // 文章列表
}
