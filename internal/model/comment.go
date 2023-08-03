package model

// Comment 评论表
type Comment struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;comment:评论ID"`
	Content   string `gorm:"type:longtext;comment:评论内容"`
	UserID    uint64 `gorm:"comment:评论用户ID"`
	ArticleID uint64 `gorm:"comment:[外键]文章ID"`
	ParentID  uint64 `gorm:"index;not null;comment:父级评论ID"`
	Floor     uint32 `gorm:"index;not null;comment:评论楼层"`
	State     uint8  `gorm:"comment:该评论状态"`

	// 预加载模型
	User User // 评论所属的用户信息，通过预加载获取

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
