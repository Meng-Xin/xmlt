package service

import (
	"context"
	"xmlt/internal/domain"
	"xmlt/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (domain.Article, error)
	Publish(ctx context.Context, article domain.Article) (domain.Article, error)
	Get(ctx context.Context, id uint64) (domain.Article, error)
}

type articleService struct {
	bRepo repository.ArticleRepo // 作者库
	cRepo repository.ArticleRepo // 线上库
}

func NewArticleService(bRepo repository.ArticleRepo, cRepo repository.ArticleRepo) ArticleService {
	return &articleService{
		bRepo: bRepo,
		cRepo: cRepo,
	}
}

// Save 在 service 层面上才会有创建或者更新的概念。repository 的职责更加单纯一点
func (a *articleService) Save(ctx context.Context, article domain.Article) (domain.Article, error) {
	if article.ID == 0 {
		id, err := a.bRepo.Create(ctx, article)
		if err != nil {
			return domain.Article{}, err
		}
		article.ID = id
		return article, nil
	}
	err := a.cRepo.Update(ctx, article)
	return article, err
}

// Publish B 端 作者进行推送文章
func (a *articleService) Publish(ctx context.Context, article domain.Article) (domain.Article, error) {
	// 先保存作者库，在推送到线上库
	art, err := a.Save(ctx, article)
	if err != nil {
		return domain.Article{}, err
	}
	// 同步给线上库，因为是不同的数据库无法使用事务，
	// 只能考虑重试 + 监控 + 告警
	_, err = a.cRepo.CreateAndCached(ctx, art)
	return art, err
}

// Get 这是 C 端查看
func (a *articleService) Get(ctx context.Context, id uint64) (domain.Article, error) {
	return a.cRepo.Get(ctx, id)
}
