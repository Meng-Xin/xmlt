package router

type RouteGroup struct {
	ArticleRouter
	UserRouter
}

var RouterGroupAll = new(RouteGroup)
