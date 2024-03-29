package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"xmlt/internal/domain"
	"xmlt/internal/model"
	"xmlt/internal/repository"
	"xmlt/utils"
)

func TestUserService(t *testing.T) {
	hashPwd := utils.MD5V("yyz", "salt", 10)
	//inputUser := domain.User{UserName: "xxn", Password: "yyz", NickName: "小仙女郁郁症"}
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) repository.UserRepo
		inputUser domain.User

		wantUser domain.User
		wantErr  error
	}{
		{
			name:      "注册成功案例",
			inputUser: domain.User{UserName: "xxn", Password: "yyz", NickName: "小仙女郁郁症"},
			mock: func(ctrl *gomock.Controller) repository.UserRepo {
				repo := repository.NewMockUserRepo(ctrl)
				repo.EXPECT().GetByUserName(gomock.Any(), "xxn").Return(model.User{}, gorm.ErrRecordNotFound)
				repo.EXPECT().GetByNickName(gomock.Any(), "小仙女郁郁症").Return(model.User{}, gorm.ErrRecordNotFound)
				repo.EXPECT().Create(gomock.Any(), domain.User{UserName: "xxn", Password: hashPwd, NickName: "小仙女郁郁症"}).Return(uint64(1), nil)
				return repo
			},

			wantUser: domain.User{ID: 1, UserName: "xxn", Password: hashPwd, NickName: "小仙女郁郁症"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 通过Mock 获得Repo构造对象
			repo := tc.mock(ctrl)
			// 下面创建Service进行测试
			svc := NewUserService(repo)
			// 下边进行 Service层接口调试
			user, err := svc.Register(context.Background(), tc.inputUser)
			// 错误一致
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			// 数据一致
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
