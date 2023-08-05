package model

// Notification 通知表
type Notification struct {
	ID           uint64 `gorm:"primaryKey;comment:消息通知ID"`
	Action       string `gorm:"size:50;comment:通知操作"`
	SubjectID    uint8  `gorm:"comment:主题ID对应#IM、CommentId、ArticleId、LikeId"`
	SubjectType  string `gorm:"size:10;comment:主题类型：IM(暂无)、COMMENT、PUSH、LIKE"`
	AcceptUserID uint64 `gorm:"index;not null;comment:接收用户"`
	SendUserID   uint64 `gorm:"index;not null;comment:发送用户"`
	State        bool   `gorm:"default:false;comment:是否查看"`
	Ctime        int64  // 创建时间，毫秒作为单位
	Utime        int64  // 更新时间，毫秒作为单位
}
