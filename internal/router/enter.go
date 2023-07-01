package router

type RouteGroup struct {
	ArticleRouter
	UserRouter
}

var RouterGroupCtrl = new(RouteGroup)
