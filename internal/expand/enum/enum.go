package enum

const (
	ArticleGetSourceOnline = 0 // 从线上库获取
	ArticleGetSourceMake   = 1 // 从制作库库获取
)

type CommentState = uint8

const (
	Normal CommentState = 0 // 正常状态
	Delete CommentState = 1 // 删除状态
)

type RangeOrder = uint8

const (
	Positive RangeOrder = 0 // 正序获取ZAdd集合内部文件
	Reverse  RangeOrder = 1 // 倒叙获取ZAdd集合的文明
)

type PageNum = int

const (
	MaxLimitNum PageNum = 20 // 最大分页获取数量
)