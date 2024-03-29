package domain

import "time"

type Article struct {
	ID           uint64 // 帖子ID
	Title        string // 帖子标题
	Content      string // 帖子内容
	CommentCount uint64 // 评论数量
	Status       uint8  // 帖子状态

	Author      uint64 // 作者
	CategoryID  uint64 // 所属板块
	NiceTopic   uint8  // 精选话题
	BrowseCount uint64 // 浏览量
	ThumbsUP    uint64 // 点赞数

	User      User
	Comments  []Comment
	UserLikes []UserLikeArticle
	Ctime     time.Time
	Utime     time.Time
}
