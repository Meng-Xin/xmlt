package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"xmlt/internal/domain"
	"xmlt/internal/shared/e"
	"xmlt/utils"
)

type UserCache interface {
	Set(ctx context.Context, user domain.User) error
	Get(ctx context.Context, uid uint64) (domain.User, error)
}

type userRedisCache struct {
	client *redis.Client
}

func (u *userRedisCache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	res, err := u.client.WithContext(ctx).Set(fmt.Sprintf("user_%d", user.ID), data, utils.GetRandTime(time.Hour, time.Minute, 30)).Result()
	if res != "OK" {
		return e.RedisInsertError
	}
	return err
}

func (u *userRedisCache) Get(ctx context.Context, uid uint64) (domain.User, error) {
	data, err := u.client.WithContext(ctx).Get(fmt.Sprintf("user_%d", uid)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(data, user)
	return user, err
}

func NewUserRedisCache(client *redis.Client) UserCache {
	return &userRedisCache{client: client}
}
