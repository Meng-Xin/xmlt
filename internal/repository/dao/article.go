package dao

import (
	"context"
	"gorm.io/gorm"
	"xmlt/internal/domain"
)

type ArticleDAO interface {
	// Create 创建一篇文章
	Create(ctx context.Context, article domain.Article) (uint64, error)
	CreateAndCached(ctx context.Context, article domain.Article) (uint64, error)
	// Update 更新一篇文章
	Update(ctx context.Context, article domain.Article) error
	// Get 方法应该负责将 Author 也一并组装起来。
	// 这里就会有一个很重要的概念，叫做延迟加载，但是 GO 是做不了的，
	// 所以只能考虑传递标记位，或者使用新方法来控制要不要把 Author 组装起来
	Get(ctx context.Context, id uint64) (domain.Article, error)
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &articleGORM{db: db}
}

type articleGORM struct {
	db *gorm.DB
}

func (a articleGORM) Create(ctx context.Context, article domain.Article) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (a articleGORM) CreateAndCached(ctx context.Context, article domain.Article) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (a articleGORM) Update(ctx context.Context, article domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (a articleGORM) Get(ctx context.Context, id uint64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
