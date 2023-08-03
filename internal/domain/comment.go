package domain

type Comment struct {
	ID        uint64 // 评论ID
	Content   string // 评论内容
	UserID    uint64 // 评论用户ID
	ArticleID uint64 // 文章ID
	ParentID  uint64 // 父级评论ID
	Floor     uint32 // 评论楼层
	State     uint8  // 该评论状态 0:正常，1：删除
	Ctime     int64  // 创建时间，毫秒作为单位
	Utime     int64  // 更新时间，毫秒作为单位
}
