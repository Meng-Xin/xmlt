package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/utils"
)

type UserLikeArticleCache interface {
	// ZAdd 点赞
	ZAdd(ctx context.Context, articleID, userID uint64) error
	// ZRem 取消点赞
	ZRem(ctx context.Context, articleID, userID uint64) error
	// ZCard 点赞数量
	ZCard(ctx context.Context, articleID uint64) (uint64, error)
	// ZRange 排序
	ZRange(ctx context.Context, articleID uint64) ([]uint64, error)
	// ZScore 返回点赞用户的文章列表ID
	ZScore(ctx context.Context, articleID, userID uint64) (float64, error)
}

func NewUserLikeArticleCache(client *redis.Client) UserLikeArticleCache {
	return &userLikeArticleCache{client: client}
}

type userLikeArticleCache struct {
	client *redis.Client
}

func (u *userLikeArticleCache) ZScore(ctx context.Context, articleID, userID uint64) (float64, error) {
	result, err := u.client.WithContext(ctx).ZScore(fmt.Sprintf(enum.UserLikeArticle, articleID), strconv.FormatInt(int64(userID), 10)).Result()
	if err == nil {
		return result, nil
	} else if err == redis.Nil {
		return 0, e.NotFoundUserLikeArticle
	}
	return 0, err
}

func (u *userLikeArticleCache) ZAdd(ctx context.Context, articleID, userID uint64) error {
	now := utils.GetTimeMilli()
	err := u.client.WithContext(ctx).ZAdd(fmt.Sprintf(enum.UserLikeArticle, articleID), redis.Z{
		Score:  float64(now),
		Member: userID,
	}).Err()
	return err
}

func (u *userLikeArticleCache) ZRem(ctx context.Context, articleID, userID uint64) error {
	// 删除当前点赞记录
	err := u.client.WithContext(ctx).ZRem(fmt.Sprintf(enum.UserLikeArticle, articleID), strconv.FormatInt(int64(userID), 10)).Err()
	return err
}

func (u *userLikeArticleCache) ZCard(ctx context.Context, articleID uint64) (uint64, error) {
	result, err := u.client.WithContext(ctx).ZCard(fmt.Sprintf(enum.UserLikeArticle, articleID)).Result()
	if err != nil {
		return 0, err
	}
	return uint64(result), err
}

func (u *userLikeArticleCache) ZRange(ctx context.Context, articleID uint64) ([]uint64, error) {
	data, err := u.client.WithContext(ctx).ZRange(fmt.Sprintf(enum.UserLikeArticle, articleID), 0, 5).Result()
	if err != nil {
		return nil, err
	}
	var userList []uint64
	for i, _ := range data {
		id, err := strconv.ParseUint(data[i], 10, 64)
		if err != nil {
			return nil, err
		}
		userList = append(userList, id)
	}
	return userList, err
}
