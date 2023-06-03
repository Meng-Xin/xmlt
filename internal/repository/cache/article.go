package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
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

func (a *articleRedisCache) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	res, err := a.client.Set(fmt.Sprintf("article_%d", article.ID), string(data), time.Hour).Result()
	if res != "OK" {
		return errors.New("Reids - 插入失败")
	}
	return err
}

func (a *articleRedisCache) Get(ctx context.Context, id uint64) (domain.Article, error) {
	// 这里之前存入的是Json 化的数据，需要反向解析绑定
	data, err := a.client.Get(fmt.Sprintf("article_%d", id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(data, art)
	return art, err
}
