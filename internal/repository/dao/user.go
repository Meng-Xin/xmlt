package dao

import (
	"context"
	"gorm.io/gorm"
	"xmlt/internal/model"
	"xmlt/utils"
)

type UserDao interface {
	// Insert 注册账号
	Insert(ctx context.Context, user model.User) (uint64, error)
	// Update 更新用户信息
	Update(ctx context.Context, user model.User) error
	// GetByID 通过ID获取用户信息
	GetByID(ctx context.Context, id uint64) (model.User, error)
	// GetByNickName 通过NickName获取用户信息
	GetByNickName(ctx context.Context, nickName string) (model.User, error)
	// GetByUserName 通过UserName获取账号信息
	GetByUserName(ctx context.Context, userName string) (model.User, error)
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userGORM{db: db}
}

type userGORM struct {
	db *gorm.DB
}

func (u *userGORM) Insert(ctx context.Context, user model.User) (uint64, error) {
	err := u.db.WithContext(ctx).Create(&user).Error
	return user.ID, err
}

func (u *userGORM) Update(ctx context.Context, user model.User) error {
	user.Utime = utils.GetTimeMilli()
	return u.db.WithContext(ctx).Model(&user).Updates(user).Error
}

func (u *userGORM) GetByID(ctx context.Context, id uint64) (model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("id=?", id).First(&user).Error
	return user, err
}

func (u *userGORM) GetByNickName(ctx context.Context, nickName string) (model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("nick_name=?", nickName).First(&user).Error
	return user, err
}

func (u *userGORM) GetByUserName(ctx context.Context, userName string) (model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("user_name=?", userName).First(&user).Error
	return user, err
}
