package e

var ErrMsg = map[int]string{
	SUCCESS:      "Ok",
	ERROR:        "内部错误",
	NotFindToken: "Token 不存在",
	TokenExpired: "Token 已过期",
	RefreshToken: "Token 无感刷新",
}

func GetMsg(code int) string {
	return ErrMsg[code]
}
