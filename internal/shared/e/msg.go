package e

var ErrMsg = map[int]string{
	SUCCESS:             "Ok",
	ERROR:               "内部错误",
	NotFindToken:        "Token 不存在",
	TokenExpired:        "Token 已过期",
	RefreshToken:        "Token 无感刷新",
	ArticleDoseNotExist: "文章不存在",
	UserExisted:         "昵称已存在",
	Registered:          "账户已被注册",
	UserOrPasswordError: "账号或密码错误",
}

func GetMsg(code int) string {
	return ErrMsg[code]
}
