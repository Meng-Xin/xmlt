package cache

import (
	"context"
	"github.com/go-redis/redis"
	"xmlt/internal/domain"
)

type ArticleCache interface {
	// Set 理论上来说，ArticleCache 也应该有自己的 Article 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id uint64) (domain.Article, error)
}

func NewArticleRedisCache(client *redis.Client) ArticleCache {
	return &articleRedisCache{
		client: client,
	}
}

type articleRedisCache struct {
	client *redis.Client
}

func (a articleRedisCache) Set(ctx context.Context, article domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (a articleRedisCache) Get(ctx context.Context, id uint64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
