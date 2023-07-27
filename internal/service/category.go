package service

import (
	"context"
	"errors"
	"xmlt/internal/domain"
	"xmlt/internal/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category domain.Category) (uint64, error)
	UpdateCategory(ctx context.Context, category domain.Category) error
	DeleteCategory(ctx context.Context, categoryID uint64) error
	GetCategoryList(ctx context.Context) ([]domain.Category, error)
}

func NewCategoryService(repo repository.CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

type categoryService struct {
	repo repository.CategoryRepo
}

func (c *categoryService) CreateCategory(ctx context.Context, category domain.Category) (uint64, error) {
	if category.ID != 0 {
		return 0, errors.New("已存在主题，无法继续创建")
	}
	return c.repo.Create(ctx, category)
}

func (c *categoryService) UpdateCategory(ctx context.Context, category domain.Category) error {
	return c.repo.Update(ctx, category)
}

func (c *categoryService) DeleteCategory(ctx context.Context, categoryID uint64) error {
	return c.repo.Delete(ctx, categoryID)
}

func (c *categoryService) GetCategoryList(ctx context.Context) ([]domain.Category, error) {
	return c.repo.GetCategoryList(ctx)
}
