package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"xmlt/internal/domain"
	"xmlt/internal/model"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/shared/e"
	"xmlt/utils"
)

type CategoryRepo interface {
	Create(ctx context.Context, category domain.Category) (uint64, error)
	Update(ctx context.Context, category domain.Category) error
	Delete(ctx context.Context, categoryID uint64) error
	GetCategoryList(ctx context.Context) ([]domain.Category, error)
}

func NewCategoryRepo(dao dao.CategoryDao, cache cache.CategoryCache) CategoryRepo {
	return &categoryRepo{dao: dao, cache: cache}
}

type categoryRepo struct {
	dao   dao.CategoryDao
	cache cache.CategoryCache
}

func (c *categoryRepo) Create(ctx context.Context, category domain.Category) (uint64, error) {
	now := utils.GetTimeMilli()
	entity := model.Category{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		ArticleCount: category.ArticleCount,
		State:        true,
		Ctime:        now,
	}
	return c.dao.Insert(ctx, entity)
}

func (c *categoryRepo) Update(ctx context.Context, category domain.Category) error {
	now := utils.GetTimeMilli()
	entity := model.Category{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		ArticleCount: category.ArticleCount,
		Utime:        now,
	}
	return c.dao.Update(ctx, entity)
}

func (c *categoryRepo) Delete(ctx context.Context, categoryID uint64) error {
	return c.dao.Delete(ctx, categoryID)
}

func (c *categoryRepo) GetCategoryList(ctx context.Context) ([]domain.Category, error) {
	cacheData, err := c.cache.Get(ctx)
	if err == nil && len(cacheData) > 0 {
		return cacheData, nil
	}
	daoData, err := c.dao.GetCategoryList(ctx)
	if err != nil {
		return nil, err
	}
	var entities []domain.Category
	for i, _ := range daoData {
		entity := daoData[i]
		entities = append(entities, c.dao2Dto(entity))
	}
	err = c.cache.Set(ctx, entities)
	if err != nil {
		log.Warning(e.RedisInsertError)
	}
	return entities, nil
}

func (c *categoryRepo) dao2Dto(entity model.Category) domain.Category {
	domainCategory := domain.Category{
		ID:           entity.ID,
		Name:         entity.Name,
		Description:  entity.Description,
		ArticleCount: entity.ArticleCount,
		Utime:        entity.Utime,
	}
	for _, daoArticle := range entity.Articles {
		domainArticle := domain.Article{
			ID:           daoArticle.ID,
			Title:        daoArticle.Title,
			Content:      daoArticle.Content,
			CommentCount: daoArticle.CommentCount,
			Status:       daoArticle.Status,
			CategoryID:   daoArticle.CategoryID,
			NiceTopic:    daoArticle.NiceTopic,
			BrowseCount:  daoArticle.BrowseCount,
			ThumbsUP:     daoArticle.ThumbsUP,
		}
		domainCategory.Articles = append(domainCategory.Articles, domainArticle)
	}
	return domainCategory
}
