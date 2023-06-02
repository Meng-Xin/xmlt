package service

import (
	"context"
	"xmlt/internal/domain"
	"xmlt/internal/repository"
)

type article struct {
	bRepo repository.ArticleRepo // 作者库
	cRepo repository.ArticleRepo // 线上库
}

func NewArticleService(bRepo repository.ArticleRepo, cRepo repository.ArticleRepo) ArticleService {
	return &article{
		bRepo: bRepo,
		cRepo: cRepo,
	}
}

// Save 在 service 层面上才会有创建或者更新的概念。repository 的职责更加单纯一点
func (a article) Save(ctx context.Context, article domain.Article) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// Publish B 端 作者进行推送文章
func (a article) Publish(ctx context.Context, article domain.Article) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// Get 这是 C 端查看
func (a article) Get(ctx context.Context, id uint64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
