package dao

import (
	"context"
	"gorm.io/gorm"
	"xmlt/internal/model"
)

type CategoryDao interface {
	Insert(ctx context.Context, category model.Category) (uint64, error)
	Update(ctx context.Context, category model.Category) error
	Delete(ctx context.Context, categoryID uint64) error
	GetCategoryList(ctx context.Context) ([]model.Category, error)
}

func NewCategoryDao(db *gorm.DB) CategoryDao {
	return &categoryGorm{db: db}
}

type categoryGorm struct {
	db *gorm.DB
}

func (c *categoryGorm) Insert(ctx context.Context, category model.Category) (uint64, error) {
	err := c.db.WithContext(ctx).Create(&category).Error
	return category.ID, err
}

func (c *categoryGorm) Update(ctx context.Context, category model.Category) error {
	return c.db.WithContext(ctx).Model(&category).Updates(category).Error
}

func (c *categoryGorm) Delete(ctx context.Context, categoryID uint64) error {
	return c.db.WithContext(ctx).Model(&model.Category{}).Where("id=?", categoryID).Update("state", "0").Error
}

func (c *categoryGorm) GetCategoryList(ctx context.Context) ([]model.Category, error) {
	var categoryList []model.Category
	err := c.db.WithContext(ctx).Select(&categoryList).Error
	return categoryList, err
}
