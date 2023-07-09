package dao

import (
	"context"
	"gorm.io/gorm"
	"xmlt/internal/domain"
	"xmlt/internal/model"
)

type CommentDao interface {
	// Insert 创建评论
	Insert(ctx context.Context, comment model.Comment) (uint64, error)
	// Update 修改评论(这里认为只有 评论状态可以修改)
	Update(ctx context.Context, comment model.Comment) error
	// GetByID 根据评论ID获取评论
	GetByID(ctx context.Context, cid uint64) (model.Comment, error)
	// GetByParentID 根据 父级评论ID 获取评论
	GetByParentID(ctx context.Context, parentID uint64, paging domain.Paging) ([]model.Comment, error)
	// GetByArticleID 根据 帖子ID 获取评论
	GetByArticleID(ctx context.Context, articleID uint64, paging domain.Paging) ([]model.Comment, error)
	// GetByUserID 根据 用户ID 获取评论
	GetByUserID(ctx context.Context, userID uint64, paging domain.Paging) ([]model.Comment, error)
}

type commentGORM struct {
	db *gorm.DB
}

func (c *commentGORM) Insert(ctx context.Context, comment model.Comment) (uint64, error) {
	err := c.db.WithContext(ctx).Create(&comment).Error
	return comment.ID, err
}

func (c *commentGORM) Update(ctx context.Context, comment model.Comment) error {
	return c.db.WithContext(ctx).Updates(&comment).Error
}

func (c *commentGORM) GetByID(ctx context.Context, cid uint64) (model.Comment, error) {
	var comment model.Comment
	err := c.db.WithContext(ctx).Where("id=?", cid).First(&comment).Error
	return comment, err
}

func (c *commentGORM) GetByParentID(ctx context.Context, parentID uint64, paging domain.Paging) ([]model.Comment, error) {
	var comments []model.Comment
	err := c.db.WithContext(ctx).Where("parent_id=?", parentID).Offset(paging.Offset).Limit(paging.Limit).Find(&comments).Error
	return comments, err
}

func (c *commentGORM) GetByArticleID(ctx context.Context, articleID uint64, paging domain.Paging) ([]model.Comment, error) {
	var comments []model.Comment
	err := c.db.WithContext(ctx).Where("article_id=?", articleID).Offset(paging.Offset).Limit(paging.Limit).Find(comments).Error
	return comments, err
}

func (c *commentGORM) GetByUserID(ctx context.Context, userID uint64, paging domain.Paging) ([]model.Comment, error) {
	var comments []model.Comment
	err := c.db.WithContext(ctx).Where("user_id=?", userID).Offset(paging.Offset).Limit(paging.Limit).Find(comments).Error
	return comments, err
}

func NewCommentDao(db *gorm.DB) CommentDao {
	return &commentGORM{db: db}
}
