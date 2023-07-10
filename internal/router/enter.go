package router

type RouteGroup struct {
	ArticleRouter
	UserRouter
	CommentRouter
}

var RouterGroupCtrl = new(RouteGroup)
