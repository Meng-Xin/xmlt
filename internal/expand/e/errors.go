package e

import "errors"

var (
	RedisInsertError     = errors.New("Reids - 插入失败")
	UserNameExistedError = errors.New("该账号已被使用")
	NickNameExistedError = errors.New("用户昵称已存在")
	UserOrPwdError       = errors.New("账号或密码不存在")
)
