package service

import (
	"context"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (uint64, error)
	Publish(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id uint64, source int) (domain.Article, error)
}

type articleService struct {
	bRepo repository.ArticleRepo // 线上库
	cRepo repository.ArticleRepo // 作者库
}

func NewArticleService(bRepo repository.ArticleRepo, cRepo repository.ArticleRepo) ArticleService {
	return &articleService{
		bRepo: bRepo,
		cRepo: cRepo,
	}
}

// Save 在 service 层面上才会有创建或者更新的概念。repository 的职责更加单纯一点
func (a *articleService) Save(ctx context.Context, article domain.Article) (uint64, error) {
	// 保存文章到制作库
	if article.ID == 0 {
		return a.cRepo.Create(ctx, article)
	}
	dataInfo, err := a.cRepo.Get(ctx, article.ID)
	if err != nil {
		return 0, err
	}
	// 文章是否存在校验
	if dataInfo.ID == 0 {
		return 0, e.UpdateArticleNotFindError
	}
	// 本次更新作者是否合法
	if dataInfo.Author != article.Author {
		return 0, e.UpdateArticleIdenticalError
	}
	// 更新制作库
	return article.ID, a.cRepo.Update(ctx, article)

}

// Publish B 端 作者进行推送文章
func (a *articleService) Publish(ctx context.Context, article domain.Article) error {
	// 先保存作者库，在推送到线上库
	id, err := a.Save(ctx, article)
	if err != nil {
		return err
	}
	article.ID = id
	// 同步给线上库，因为是不同的数据库无法使用事务，
	// 只能考虑重试 + 监控 + 告警
	_, err = a.bRepo.CreateAndCached(ctx, article)
	return err
}

// Get Source：0-线上库 1-制作库
func (a *articleService) Get(ctx context.Context, id uint64, source int) (domain.Article, error) {
	if source == 0 {
		return a.bRepo.Get(ctx, id)
	}
	return a.cRepo.Get(ctx, id)
}
