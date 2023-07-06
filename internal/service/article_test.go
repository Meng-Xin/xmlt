package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"xmlt/internal/domain"
	repository "xmlt/internal/repository"
)

func TestService_Save(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.ArticleRepo
		inputArt domain.Article

		wantArt domain.Article
		wantErr error
	}{
		{
			name: "创建",
			inputArt: domain.Article{Author: 123,
				Title: "这是标题", Content: "这是内容"},
			mock: func(ctrl *gomock.Controller) repository.ArticleRepo {
				repo := repository.NewMockArticleRepo(ctrl)
				repo.EXPECT().
					Create(gomock.Any(), domain.Article{Author: 123,
						Title: "这是标题", Content: "这是内容"}).
					Return(uint64(1), nil)
				return repo
			},
			wantArt: domain.Article{ID: 1, Title: "这是标题", Content: "这是内容",
				Author: 123},
		},
		// 你可以在这里试试写一下更新的测试用例
		//{
		//	name:     "更新",
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			// 我们 Save 是只用 ToB 的
			svc := NewArticleService(repo, nil)
			// 这里我们直接写死了 ctx，因为不打算测试超时或者取消之类的
			art, err := svc.Save(context.Background(), tc.inputArt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantArt, art)
		})
	}
}
