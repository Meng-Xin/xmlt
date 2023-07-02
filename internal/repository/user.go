package repository

import (
	"context"
	"errors"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/model"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/utils"
)

type UserRepo interface {
	// Create 注册
	Create(ctx context.Context, user domain.User) (uint64, error)
	// Login 登录
	Login(ctx context.Context, userName string) (domain.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user domain.User) error
	// Get 获取用户信息
	Get(ctx context.Context, id uint64) (domain.User, error)
	// GetByNickName 通过NickName获取用户信息
	GetByNickName(ctx context.Context, nickName string) (model.User, error)
	// GetByUserName 通过UserName获取用户信息
	GetByUserName(ctx context.Context, userName string) (model.User, error)
}

type userRepo struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func (u *userRepo) Create(ctx context.Context, user domain.User) (uint64, error) {
	now := utils.GetTimeMilli()
	// 创建用户
	if user.ID != 0 {
		return 0, errors.New("用户已存在无法创建")
	}
	entity := model.User{
		ID:       user.ID,
		UserName: user.UserName,
		Password: user.Password,
		NickName: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		Ctime:    now,
		Utime:    now,
	}
	return u.dao.Insert(ctx, entity)
}

func (u *userRepo) Login(ctx context.Context, userName string) (domain.User, error) {
	// 因为这里是登录，最大访问量也是用户自己对自己访问，所以走cache
	user, err := u.dao.GetByUserName(ctx, userName)
	if err != nil {
		return domain.User{}, err
	}
	entity := domain.User{
		ID:       user.ID,
		UserName: user.UserName,
		Password: user.Password,
		NickName: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	}
	return entity, nil
}

func (u *userRepo) Update(ctx context.Context, user domain.User) error {
	now := utils.GetTimeMilli()
	if user.ID == 0 {
		return errors.New("更新对象不存在！")
	}
	// 先查询拿到原有数据
	entity, err := u.dao.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	// 再对源数据部分字段更新
	entity.NickName = user.NickName
	entity.Avatar = user.Avatar
	entity.Utime = now
	return u.dao.Update(ctx, entity)
}

func (u *userRepo) Get(ctx context.Context, id uint64) (domain.User, error) {
	// 先查看缓存
	res, err := u.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}
	// 缓存不命中,执行数据库查询,再保存到缓存
	entity, err := u.dao.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:       entity.ID,
		UserName: entity.UserName,
		Password: entity.Password,
		NickName: entity.NickName,
		Email:    entity.Email,
		Phone:    entity.Phone,
		Avatar:   entity.Avatar,
	}
	// 更新缓存
	err = u.cache.Set(ctx, user)
	if err != nil {
		// 数据库更新成功，但是缓存更新失败，记录日志。
		global.Log.Error(err)
	}
	return user, nil
}

func (u *userRepo) GetByNickName(ctx context.Context, nickName string) (model.User, error) {
	return u.dao.GetByNickName(ctx, nickName)
}
func (u *userRepo) GetByUserName(ctx context.Context, userName string) (model.User, error) {
	return u.dao.GetByUserName(ctx, userName)
}
func NewUserRepo(dao dao.UserDao, cache cache.UserCache) UserRepo {
	return &userRepo{dao: dao, cache: cache}
}
