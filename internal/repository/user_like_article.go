package repository

import (
	"context"
	"xmlt/internal/expand/e"
	"xmlt/internal/repository/cache"
)

type UserLikeArticleRepo interface {
	Update(ctx context.Context, artID, uID uint64) error
	GetLikeState(ctx context.Context, artID, uID uint64) bool
	GetLikeNum(ctx context.Context, artID uint64) (uint64, error)
}

func NewUserLikeArticleRepo(cache cache.UserLikeArticleCache) UserLikeArticleRepo {
	return &userLikeArticleRepo{cache: cache}
}

type userLikeArticleRepo struct {
	cache cache.UserLikeArticleCache
}

func (u *userLikeArticleRepo) GetLikeState(ctx context.Context, artID, uID uint64) bool {
	score, err := u.cache.ZScore(ctx, artID, uID)
	if err == nil && score != 0 {
		return true
	}
	return false
}

func (u *userLikeArticleRepo) Update(ctx context.Context, artID, uID uint64) error {
	_, err := u.cache.ZScore(ctx, artID, uID)
	if err == nil {
		// 存在就删除点赞信息
		err = u.cache.ZRem(ctx, artID, uID)
		if err != nil {
			return err
		}
	} else if err == e.NotFoundUserLikeArticle {
		// 不存在就创建
		err = u.cache.ZAdd(ctx, artID, uID)
		if err != nil {
			return err
		}
	}
	return err
}

func (u *userLikeArticleRepo) GetLikeNum(ctx context.Context, artID uint64) (uint64, error) {
	return u.cache.ZCard(ctx, artID)
}
