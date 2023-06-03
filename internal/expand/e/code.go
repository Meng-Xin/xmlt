package e

const (
	SUCCESS             = 200   // OK
	ERROR               = 400   // 内部能错误
	NotFindToken        = 10010 // Token 不存在
	TokenExpired        = 10011 // Token 已过期
	RefreshToken        = 10012 // Token 无感刷新
	ArticleDoseNotExist = 10013 // 文章不存在
)
