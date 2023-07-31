package service

import (
	"context"
	"xmlt/internal/domain"
	"xmlt/internal/expand/e"
	"xmlt/internal/expand/enum"
	"xmlt/internal/repository"
	"xmlt/utils"
)

type CommentService interface {
	// AddComment 新增评论
	AddComment(ctx context.Context, comment domain.Comment) error
	// DeleteComment 删除评论 (修改评论状态)
	DeleteComment(ctx context.Context, comment domain.Comment) error
	// GetArtComment 获取帖子评论
	GetArtComment(ctx context.Context, articleID uint64, paging *domain.Page, by domain.RangeBy) ([]domain.Comment, error)
	// GetReplyComment 获取回复评论
	GetReplyComment(ctx context.Context, commentID uint64) (domain.Comment, error)
	// GetUserComment 获取用户评论
	GetUserComment(ctx context.Context, userID uint64, paging *domain.Page) ([]domain.Comment, error)
}

type commentService struct {
	repo repository.CommentRepo
}

func (c *commentService) AddComment(ctx context.Context, comment domain.Comment) error {
	// 新增评论 TODO 待完善，需要走消息队保证楼层并发问题
	return c.repo.CreateAndCached(ctx, comment)
}

func (c *commentService) DeleteComment(ctx context.Context, comment domain.Comment) error {
	data, err := c.repo.Get(ctx, comment.ID)
	// 错误校验
	if err != nil {
		return err
	}
	// 身份校验
	if data.UserID != comment.UserID {
		return e.UpdateCommentIdentityError
	}
	// 删除评论(修改评论状态)
	return c.repo.Update(ctx, comment)
}

func (c *commentService) GetArtComment(ctx context.Context, articleID uint64, paging *domain.Page, by domain.RangeBy) ([]domain.Comment, error) {
	// 参数合法校验
	//if paging.Limit == 0 {
	//	return nil, e.PagingIsNullError
	//}
	if by.Start == 0 && by.Stop == 0 {
		return nil, e.PagingIsNullError
	}
	//if paging.Offset > enum.MaxLimitNum {
	//	return nil, e.PagingMaxLimitError
	//}
	if int(utils.Abs(by.Start-by.Stop)) > enum.MaxPageSize {
		return nil, e.PagingMaxLimitError
	}
	// 获取评论信息
	return c.repo.GetByArticleID(ctx, articleID, paging, by)
}

func (c *commentService) GetReplyComment(ctx context.Context, commentID uint64) (domain.Comment, error) {
	return c.repo.Get(ctx, commentID)
}

func (c *commentService) GetUserComment(ctx context.Context, userID uint64, paging *domain.Page) ([]domain.Comment, error) {
	return c.repo.GetByUserID(ctx, userID, paging)
}

func NewCommentService(repo repository.CommentRepo) CommentService {
	return &commentService{repo: repo}
}
