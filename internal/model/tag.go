package model

// Tag 标签表
type Tag struct {
	ID       uint16    `gorm:"primaryKey;comment:所属标签ID"`
	Name     string    `gorm:"size:20;comment:标签名称"`
	Articles []Article `gorm:"many2many:article_tag"` // Tag : Article -> N : N

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
