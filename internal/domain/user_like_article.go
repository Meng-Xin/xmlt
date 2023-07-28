package domain

type UserLikeArticle struct {
	ID        uint64
	UserID    uint64
	ArticleID uint64
	LikeState bool

	Ctime uint64
	Utime uint64
}
