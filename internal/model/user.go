package model

type User struct {
	ID       uint64 `gorm:"primaryKey,autoIncrement"`
	UserName string `gorm:"size:50"`  // 账号
	Password string `gorm:"size:100"` // 密码
	NickName string `gorm:"size:30"`  // 昵称
	Email    string `gorm:"size:100"` // 邮箱
	Phone    string `gorm:"size:11"`  // 手机号
	Avatar   string // 头像

	Ctime int64 // 创建时间，毫秒作为单位
	Utime int64 // 更新时间，毫秒作为单位
}
