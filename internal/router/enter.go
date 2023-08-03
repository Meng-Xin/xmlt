package router

type RouteGroup struct {
	UserRouter
	ArticleRouter
	UserLikeArticleRouter
	CommentRouter
	CategoryRouter
	HomeRouter
}

var RouterGroupCtrl = new(RouteGroup)
