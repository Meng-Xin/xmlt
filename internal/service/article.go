package service

import (
	"context"
	"errors"
	"xmlt/internal/domain"
	"xmlt/internal/repository"
	"xmlt/internal/shared"
	"xmlt/internal/shared/e"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (uint64, error)
	Publish(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id uint64, source int) (domain.Article, error)
	GetCategoryArticles(ctx context.Context, categoryID uint64, paging *shared.Page) ([]domain.Article, error)
}

type articleService struct {
	cRepo repository.ArticleRepo // 作者库
}

func NewArticleService(cRepo repository.ArticleRepo) ArticleService {
	return &articleService{
		cRepo: cRepo,
	}
}

func (a *articleService) GetCategoryArticles(ctx context.Context, categoryID uint64, paging *shared.Page) ([]domain.Article, error) {
	if categoryID == 0 {
		return nil, errors.New("不存在该主题")
	}
	return a.cRepo.GetArticlesByCategoryID(ctx, categoryID, paging)
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

// Publish 作者进行推送文章
func (a *articleService) Publish(ctx context.Context, article domain.Article) error {
	// 先保存作者库，在推送到线上库
	id, err := a.Save(ctx, article)
	if err != nil {
		return err
	}
	article.ID = id
	// TODO 把这条消息推送到待审核文章
	return err
}

// Get Source：
func (a *articleService) Get(ctx context.Context, id uint64, source int) (domain.Article, error) {
	return a.cRepo.Get(ctx, id)
}
