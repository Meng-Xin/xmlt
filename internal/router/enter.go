package router

type RouteGroup struct {
	UserRouter
	ArticleRouter
	UserLikeArticleRouter
	CommentRouter
	CategoryRouter
}

var RouterGroupCtrl = new(RouteGroup)
