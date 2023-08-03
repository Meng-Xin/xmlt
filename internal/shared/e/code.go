package e

const (
	SUCCESS             = 200   // OK
	ERROR               = 400   // 内部能错误
	NotFindToken        = 10010 // Token 不存在
	TokenExpired        = 10011 // Token 已过期
	RefreshToken        = 10012 // Token 无感刷新
	ArticleDoseNotExist = 10013 // 文章不存在
	UserExisted         = 10014 // 用户已存在
	Registered          = 10015 // 账户已被注册
	UserOrPasswordError = 10016 // 账号或密码错误
)
