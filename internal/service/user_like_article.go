package service

import (
	"context"
	"xmlt/internal/repository"
)

type UserLikeArticleService interface {
	UpdateLike(ctx context.Context, artID, uID uint64) error
	GetLikeState(ctx context.Context, artID, uID uint64) bool
	GetLikeNum(ctx context.Context, artID uint64) (uint64, error)
}

func NewUserLikeArticleService(repo repository.UserLikeArticleRepo) UserLikeArticleService {
	return &userLikeArticleService{repo: repo}
}

type userLikeArticleService struct {
	repo repository.UserLikeArticleRepo
}

func (u *userLikeArticleService) GetLikeState(ctx context.Context, artID, uID uint64) bool {
	return u.repo.GetLikeState(ctx, artID, uID)
}

func (u *userLikeArticleService) UpdateLike(ctx context.Context, artID, uID uint64) error {
	return u.repo.Update(ctx, artID, uID)
}

func (u *userLikeArticleService) GetLikeNum(ctx context.Context, artID uint64) (uint64, error) {
	return u.repo.GetLikeNum(ctx, artID)
}
