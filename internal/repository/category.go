package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/model"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
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
		entities = append(entities, domain.Category(entity))
	}
	err = c.cache.Set(ctx, entities)
	if err != nil {
		log.Warning(e.RedisInsertError)
	}
	return entities, nil
}
