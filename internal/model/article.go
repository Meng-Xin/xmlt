package model

// Article 直接对应到表结构
type Article struct {
	ID      uint64 `gorm:"primaryKey,autoIncrement"`
	Title   string `form:"title"`
	Content string `form:"content"`
	// 作者 ID
	Author uint64 `gorm:"index,not null"`
	// 创建时间，毫秒作为单位
	Ctime int64
	// 更新时间，毫秒作为单位
	Utime int64
}
