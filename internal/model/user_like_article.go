package model

type UserLikeArticle struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;not null;comment:用户点赞表"`
	UserID    uint64 `gorm:"comment:点赞用户ID"`
	ArticleID uint64 `gorm:"comment:点赞文章ID"`
	LikeState bool   `gorm:"comment:点赞状态# 0禁用,1启用"`

	// 预加载模型
	User    User
	Article Article

	Ctime uint64 `gorm:"comment:创建时间"`
	Utime uint64 `gorm:"comment:修改时间"`
}
