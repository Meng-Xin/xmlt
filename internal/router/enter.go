package router

type RouteGroup struct {
	ArticleRouter
	UserRouter
	CommentRouter
	CategoryRouter
}

var RouterGroupCtrl = new(RouteGroup)
