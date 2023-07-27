package v1

import (
	"github.com/gin-gonic/gin"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/public"
	"xmlt/internal/service"
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
