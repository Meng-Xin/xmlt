package utils

func ContainsCompared[T Ordered](src []T, dest T) bool {
	for _, item := range src {
		if item == dest {
			return true
		}
	}
	return false
}

// DifferenceCompared 取前者src与后者dest两个字符串列表的差集
func DifferenceCompared[T Ordered](src []T, dest []T) []T {
	res := make([]T, 0)
	for _, item := range src {
		if !ContainsCompared(dest, item) {
			res = append(res, item)
		}
	}
	return res
}

// IntersectionCompared 取两个字符串列表的交集
func IntersectionCompared[T Ordered](src []T, dest []T) []T {
	res := make([]T, 0)
	for _, item := range src {
		if ContainsCompared(dest, item) {
			res = append(res, item)
		}
	}
	return res
}

// UnionComPared 取两个字符串列表的并集
func UnionComPared[T Ordered](src []T, dest []T) []T {
	res := make([]T, 0)
	res = append(res, src...)
	for _, item := range dest {
		if !ContainsCompared(res, item) {
			res = append(res, item)
		}
	}
	return res
}
