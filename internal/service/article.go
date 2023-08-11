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
	Publish(ctx context.Context, artId, uid uint64) (uint64, error)
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
func (a *articleService) Publish(ctx context.Context, artId, uid uint64) (uint64, error) {
	// TODO 把本文章的状态更改为待审核
	// 校验发布用户和文章用户是否一致
	daoArt, err := a.cRepo.Get(ctx, artId)
	if err != nil {
		return 0, err
	}
	if daoArt.Author != uid {
		return 0, errors.New("发布失败：发布用户和文章作者Id不一致！")
	}
	publishId, err := a.cRepo.UpdateToPublish(ctx, artId)
	if err != nil {
		return 0, err
	}
	return publishId, nil
}

// Get Source：
func (a *articleService) Get(ctx context.Context, id uint64, source int) (domain.Article, error) {
	return a.cRepo.Get(ctx, id)
}
