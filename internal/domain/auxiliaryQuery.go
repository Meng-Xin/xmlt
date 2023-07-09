package domain

// Paging 分页查询 辅助结构
type Paging struct {
	Offset int // 偏移量
	Limit  int // 每页限制数量
}

// RangeBy ZSet 中 Score 的范围
type RangeBy struct {
	Start int64
	Stop  int64
	Order uint8 // 0 从小到大，1 从大到小
}
