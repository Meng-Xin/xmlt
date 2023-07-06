package model

// Notification 通知表
type Notification struct {
	ID           uint64 `gorm:"primaryKey;comment:消息通知ID"`
	Action       string `gorm:"size:50;comment:通知操作"`
	SubjectID    uint8  `gorm:"comment:主题ID"`
	AcceptUserID uint64 `gorm:"index;not null;comment:接收用户"`
	SendUserID   uint64 `gorm:"index;not null;comment:发送用户"`

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
