package service

import (
	"context"
	"xmlt/internal/domain"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (domain.Article, error)
	Publish(ctx context.Context, article domain.Article) (domain.Article, error)
	Get(ctx context.Context, id uint64) (domain.Article, error)
}
