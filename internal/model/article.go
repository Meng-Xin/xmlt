package model

// Article 直接对应到表结构
type Article struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement;comment:帖子ID"`
	Title        string `gorm:"size:50;comment:帖子标题"`
	Content      string `gorm:"type:longtext;comment:帖子内容"`
	CommentCount uint64 `gorm:"comment:评论总数"`
	Status       uint8  `gorm:"comment:帖子状态 0:审核、1:正常、2:删除"`
	UserID       uint64 `gorm:"comment:作者ID"`
	CategoryID   uint64 `gorm:"comment:所属板块ID"`
	NiceTopic    uint8  `gorm:"comment:精选话题"`
	BrowseCount  uint64 `gorm:"comment:浏览量"`
	ThumbsUP     uint64 `gorm:"comment:点赞数"`

	Comments  []Comment         `gorm:"foreignKey:ArticleID;references:ID;"` // Article : Comment -> 1:N
	UserLikes []UserLikeArticle `gorm:"foreignKey:ArticleID;references:ID;"` // Article : Comment -> 1:N
	Tags      []Tag             `gorm:"many2many:article_tag"`               // Tag : Article -> N:N

	// 预加载模型
	User User

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
