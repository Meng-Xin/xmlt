package model

// Thumbs 点赞表
type Thumbs struct {
	ID        uint64 `gorm:"primaryKey;comment:点赞表ID"`
	UserID    uint64 `gorm:"index;not null;comment:点赞用户"`
	ArticleID uint64 `gorm:"index;not null;comment:点赞文章"`

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
