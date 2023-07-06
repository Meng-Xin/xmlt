package model

// Category 主题表
type Category struct {
	ID           uint16 `gorm:"primaryKey;autoIncrement;comment:板块ID"`
	Name         string `gorm:"size:30;comment:板块名称"`
	Description  string `gorm:"size:200;comment:板块描述"`
	ArticleCount uint64 `gorm:"comment:板块文章数量"`

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
