package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/utils"
)

type CategoryCache interface {
	Get(ctx context.Context) ([]domain.Category, error)
	Set(ctx context.Context, category []domain.Category) error
}

func NewCategoryCache(client *redis.Client) CategoryCache {
	return &categoryCache{client: client}
}

type categoryCache struct {
	client *redis.Client
}

func (c *categoryCache) Get(ctx context.Context) ([]domain.Category, error) {
	data, err := c.client.WithContext(ctx).Get(enum.AllCategory).Bytes()
	if err != nil {
		return []domain.Category{}, err
	}
	var res []domain.Category
	err = json.Unmarshal(data, &res)
	return res, err
}

func (c *categoryCache) Set(ctx context.Context, category []domain.Category) error {
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	result, err := c.client.WithContext(ctx).Set(enum.AllCategory, data, utils.GetRandTime(time.Hour*72, time.Minute, 45)).Result()
	if result != "OK" {
		return e.RedisInsertError
	}
	return err
}
