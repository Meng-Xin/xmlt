package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
)

// HomeController Home页面横跨了多个领域，因此它是 Service 层级的调动者
type HomeController struct {
	UserService            service.UserService
	ArticleService         service.ArticleService
	CategoryService        service.CategoryService
	UserLikeArticleService service.UserLikeArticleService
}

func NewHomeController(userService service.UserService,
	artService service.ArticleService,
	categoryService service.CategoryService,
	userLikeService service.UserLikeArticleService) *HomeController {
	return &HomeController{
		UserService:            userService,
		ArticleService:         artService,
		CategoryService:        categoryService,
		UserLikeArticleService: userLikeService,
	}
}

func (h *HomeController) GetHomeInfo(ctx *gin.Context) {
	code := e.SUCCESS
	// 获取主题列表
	var firstCate domain.Category
	categoryList, err := h.CategoryService.GetCategoryList(ctx)
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Warn(err.Error())
		return
	}
	if len(categoryList) >= 0 {
		firstCate = categoryList[0]
	}
	var artsVo []articleHome
	// 获取首页文章列表
	articles, err := h.ArticleService.GetCategoryArticles(ctx, firstCate.ID, domain.NewPage(1, 1))
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Warn(err.Error())
		return
	}
	for i, _ := range articles {
		var userVo userHome
		// 获取文章点赞数
		likeNum, err := h.UserLikeArticleService.GetLikeNum(ctx, articles[i].ID)
		if err != nil {
			code = e.ERROR
			ctx.JSON(http.StatusOK, public.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			global.Log.Warn(err.Error())
			continue
		}
		// 组装到articles内部
		articles[i].ThumbsUP = likeNum
		// 获取用户信息并组装 1.获取用户头像|名称
		user, err := h.UserService.Get(ctx, articles[i].Author)
		if err != nil {
			return
		}
		userVo.init(user)
		art := articleHome{}
		art.init(articles[i], userVo)
		artsVo = append(artsVo, art)
	}
	homeVo := HomeVo{}
	homeVo.categories.init(firstCate, artsVo)
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   homeVo,
	})
}

type HomeVo struct {
	categories categoryHome
}

type categoryHome struct {
	ID          uint64
	Name        string
	Description string
	articles    []articleHome
}

func (c *categoryHome) init(gory domain.Category, arts []articleHome) {
	c.ID = gory.ID
	c.Name = gory.Name
	c.Description = gory.Description
	c.articles = arts
}

type articleHome struct {
	ID         uint64
	Title      string
	Author     userHome
	NiceTopic  uint8
	BrowseCoun uint64
	Ctime      string
}

func (a *articleHome) init(article domain.Article, auth userHome) {
	a.ID = article.ID
	a.Title = article.Title
	a.Author = auth
	a.NiceTopic = article.NiceTopic
	a.BrowseCoun = article.BrowseCount
	a.Ctime = article.Ctime.String()
}

type userHome struct {
	ID     uint64
	Name   string
	Avatar string
}

func (u *userHome) init(user domain.User) {
	u.ID = user.ID
	u.Name = user.UserName
	u.Avatar = user.Avatar
}
