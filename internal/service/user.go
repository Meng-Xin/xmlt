package service

import (
	"context"
	"errors"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/repository"
	"xmlt/utils"
)

type UserService interface {
	// Register 注册
	Register(ctx context.Context, user domain.User) (domain.User, error)
	// Login 登录
	Login(ctx context.Context, user domain.User) (domain.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user domain.User) (domain.User, error)
	// Get 获取用户信息
	Get(ctx context.Context, id uint64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}
func (u *userService) Register(ctx context.Context, user domain.User) (domain.User, error) {
	// 账户已存在校验
	ok := u.CheckUserName(ctx, user.UserName)
	if !ok {
		return domain.User{}, e.UserNameExistedError
	}
	// 昵称重复校验
	ok = u.CheckNickName(ctx, user.NickName)
	if !ok {
		return domain.User{}, e.NickNameExistedError
	}
	// 对用户密码加盐Hash
	hashPwd := utils.MD5V(user.Password, "salt", 10)
	user.Password = hashPwd
	// 创建用户
	uid, err := u.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = uid
	return user, nil
}

func (u *userService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	// 用户登录
	data, err := u.repo.Login(ctx, user.UserName)
	if err != nil {
		return domain.User{}, err
	}
	// Hash比较
	hashPwdLogin := utils.MD5V(user.Password, "salt", 10)
	if user.UserName != data.UserName || hashPwdLogin != data.Password {
		return domain.User{}, e.UserOrPwdError
	}
	// 登陆成功，给用户颁发双Token
	aToken, rToken, err := utils.GenToken(data.ID, data.NickName)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = data.ID
	user.NickName = data.NickName
	user.Token = aToken + "," + rToken
	return user, nil
}

func (u *userService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	// 获取用户信息
	if user.ID == 0 {
		return domain.User{}, errors.New("用户不存在")
	}
	// 更新用户信息
	err := u.repo.Update(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userService) Get(ctx context.Context, id uint64) (domain.User, error) {
	// 获取用户信息
	if id == 0 {
		return domain.User{}, errors.New("用户不存在")
	}
	user, err := u.repo.Get(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userService) CheckNickName(ctx context.Context, nickName string) bool {
	user, _ := u.repo.GetByNickName(ctx, nickName)
	if user.ID != 0 {
		return false
	}
	//global.Log.Info("空查询测试：" + err.Error())
	return true
}
func (u *userService) CheckUserName(ctx context.Context, userName string) bool {
	user, _ := u.repo.GetByUserName(ctx, userName)
	if user.ID != 0 {
		return false
	}
	//global.Log.Info("空查询测试：" + err.Error())
	return true
}
