package domain

import "time"

type User struct {
	ID       uint64
	UserName string // 账号
	Password string // 密码
	NickName string // 昵称
	Email    string // 邮箱
	Phone    string // 手机号
	Avatar   string // 头像

	Token string
	Ctime time.Time
	Utime time.Time
}
