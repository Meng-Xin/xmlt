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
	"xmlt/internal/shared/public"
)

type CategoryController struct {
	service service.CategoryService
}

func NewCategoryController(service service.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	code := e.SUCCESS
	var vo category
	err := ctx.Bind(&vo)
	if err != nil {
		code := e.ERROR
		ctx.JSON(code, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Info(err.Error())
		return
	}
	categoryID, err := c.service.CreateCategory(ctx.Request.Context(), domain.Category{
		Name:        vo.Name,
		Description: vo.Description,
		State:       true,
	})
	if err != nil {
		code = e.ERROR
		ctx.JSON(code, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		global.Log.Info(err.Error())
		return
	}
	ctx.JSON(code, public.Response{
		Status: code,
		Data:   categoryID,
		Msg:    "创建成功",
	})
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {

}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {

}

func (c *CategoryController) GetCategoryList(ctx *gin.Context) {
	code := e.SUCCESS
	cateId, _ := strconv.ParseUint(ctx.Query("category_id"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	newPage := shared.NewPage(page, pageSize)
	cateArts, err := c.service.GetArticlesByCateId(ctx.Request.Context(), cateId, newPage)
	if err != nil {
		code = e.ERROR
		ctx.JSON(http.StatusOK, public.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		})
		return
	}
	ctx.JSON(http.StatusOK, public.Response{
		Status: code,
		Data:   cateArts.Articles,
	})

}

type category struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArticleCount uint64 `json:"article_count"`
	State        bool   `json:"state"`
	Ctime        int64  `json:"ctime"`
	Utime        int64  `json:"utime"`
}

type CategoryInfo struct {
	Name        string
	Description string
	Articles    ArticlesRes
}

type CreateCategoryReq struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" from:"description"`
}

type UpdateCategoryReq struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" from:"description"`
}

type CreateCategoryRes struct {
	CategoryId uint64
}

func (c *CreateCategoryRes) init(cateId uint64) {
	c.CategoryId = cateId
}

type UpdateCategoryRes struct {
	CategoryId uint64
}

func (u *UpdateCategoryRes) init(cateId uint64) {
	u.CategoryId = cateId
}

type DeleteCategoryRes struct {
	CategoryId uint64
}

func (d *DeleteCategoryRes) init(cateId uint64) {
	d.CategoryId = cateId
}

type GetCategoriesRes struct {
	CateList []CategoryInfo
}

func (g *GetCategoriesRes) init(cates []domain.Category) {
	for _, cate := range cates {
		var cateInfo CategoryInfo
		cateInfo.Name = cate.Name
		cateInfo.Description = cate.Description
		cateInfo.Articles.init(cate.Articles, nil)
		g.CateList = append(g.CateList)
	}
}
