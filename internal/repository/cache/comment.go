package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/utils"
)

type CommentCache interface {
	Set(ctx context.Context, comment domain.Comment) error
	Get(ctx context.Context, cid uint64) (domain.Comment, error)
	ZAdd(ctx context.Context, comment ...domain.Comment) error
	ZGet(ctx context.Context, cid uint64, by domain.RangeBy) ([]domain.Comment, error)
}

type commentRedisCache struct {
	client *redis.Client
}

func (c *commentRedisCache) ZAdd(ctx context.Context, comments ...domain.Comment) error {
	var err error
	for i, _ := range comments {
		comment := comments[i]
		data, err := json.Marshal(comment)
		if err != nil {
			return err
		}
		// 使用 ZSet 存储 某篇文章 的评论
		_, err = c.client.WithContext(ctx).ZAdd(fmt.Sprintf("article_comments_%d", comment.ArticleID), redis.Z{
			Score:  float64(comment.Floor),
			Member: data,
		}).Result()
		if err != nil {
			break
		}
	}
	return err
}

func (c *commentRedisCache) ZGet(ctx context.Context, articleID uint64, by domain.RangeBy) ([]domain.Comment, error) {
	var comments []domain.Comment
	var result []string
	var err error
	if by.Order == enum.Positive {
		result, err = c.client.WithContext(ctx).ZRange(
			fmt.Sprintf("article_comments_%d", articleID), by.Start, by.Stop,
		).Result()
		if len(result) == 0 {
			return nil, errors.New("空数据")
		}
		if err != nil {
			return comments, err
		}
	} else {
		result, err = c.client.WithContext(ctx).ZRevRange(
			fmt.Sprintf("article_comments_%d", articleID), by.Start, by.Stop,
		).Result()
		if len(result) == 0 {
			return nil, errors.New("空数据")
		}
		if err != nil {
			return comments, err
		}
	}

	for i, _ := range result {
		// 获取对应数据
		data := result[i]
		// 构建对象
		entity := domain.Comment{}
		err = json.Unmarshal([]byte(data), &entity)
		if err != nil {
			return comments, err
		}
		comments = append(comments, entity)
	}
	return comments, nil
}

func (c *commentRedisCache) Set(ctx context.Context, comment domain.Comment) error {
	data, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	res, err := c.client.WithContext(ctx).Set(fmt.Sprintf("comment_%d", comment.ID), data, utils.GetRandTime(24*time.Hour, time.Minute, 30)).Result()
	if res != "OK" {
		return e.RedisInsertError
	}
	return err
}

func (c *commentRedisCache) Get(ctx context.Context, cid uint64) (domain.Comment, error) {
	var comment domain.Comment
	// 从Redis获取
	bytes, err := c.client.WithContext(ctx).Get(fmt.Sprintf("comment_%d", cid)).Bytes()
	if err != nil {
		return domain.Comment{}, err
	}
	// 绑定并返回结果
	return comment, json.Unmarshal(bytes, &comment)
}

func NewCommentRedisCache(client *redis.Client) CommentCache {
	return &commentRedisCache{client: client}
}
