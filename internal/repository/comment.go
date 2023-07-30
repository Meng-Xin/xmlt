package repository

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"xmlt/global"
	"xmlt/internal/domain"
	"xmlt/internal/expand/enum"
	"xmlt/internal/model"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/utils"
)

type CommentRepo interface {
	// CreateAndCached 创建评论,并写入缓存
	CreateAndCached(ctx context.Context, comment domain.Comment) error
	// Update 更新评论状态 只能更新自己的评论状态
	Update(ctx context.Context, comment domain.Comment) error
	// Get 获取单条评论信息
	Get(ctx context.Context, id uint64) (domain.Comment, error)
	// GetByArticleID 根据 帖子ID 获取评论
	GetByArticleID(ctx context.Context, id uint64, paging domain.Page, by domain.RangeBy) ([]domain.Comment, error)
	// GetByUserID 根据 用户ID 获取评论
	GetByUserID(ctx context.Context, id uint64, paging domain.Page) ([]domain.Comment, error)
	// GetLatestFloor 获取最新楼层
	GetLatestFloor(ctx context.Context, articleID uint64) (uint32, error)
	// ConsumerMQ
	ConsumerMQ(ctx context.Context) error
}

type commentRepo struct {
	dao   dao.CommentDao
	cache cache.CommentCache
}

func (c *commentRepo) ConsumerMQ(ctx context.Context) error {
	// 从消息队列获取
	msgs, err := global.RabbitMQ.Ch.Consume(global.RabbitMQ.QueueName, "评论消费", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() error {
		for msg := range msgs {
			endData := model.Comment{}
			err = json.Unmarshal(msg.Body, &endData)
			if err != nil {
				return err
			}
			floor, err := c.GetLatestFloor(ctx, endData.ArticleID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			endData.Floor = floor + 1
			_, err = c.dao.Insert(ctx, endData)
			if err != nil {
				return err
			}
			// 写入缓存
			err = c.cache.Set(ctx, domain.Comment(endData))
			err = c.cache.ZAdd(ctx, domain.Comment(endData))
			if err != nil {
				global.Log.Warn(err.Error())
			}
		}
		return err
	}()
	return err
}

func (c *commentRepo) GetLatestFloor(ctx context.Context, articleID uint64) (uint32, error) {
	return c.dao.GetLatestFloorByArticleID(ctx, articleID)
}

func (c *commentRepo) CreateAndCached(ctx context.Context, comment domain.Comment) error {
	comment.Ctime = utils.GetTimeMilli()
	entity := model.Comment{
		Content:   comment.Content,
		UserID:    comment.UserID,
		ArticleID: comment.ArticleID,
		ParentID:  comment.ParentID,
		Floor:     comment.Floor,
		Ctime:     comment.Ctime,
	}
	// 存入消息队列
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	err = global.RabbitMQ.PublishOnQueue(ctx, data)
	if err != nil {
		return err
	}

	return err
}

func (c *commentRepo) Update(ctx context.Context, comment domain.Comment) error {
	comment.Utime = utils.GetTimeMilli()
	entity := model.Comment{
		ID:     comment.ID,
		UserID: comment.UserID,
		State:  enum.Delete,
		Utime:  comment.Utime,
	}
	return c.dao.Update(ctx, entity)
}

func (c *commentRepo) Get(ctx context.Context, id uint64) (domain.Comment, error) {
	// 先看缓存
	cacheComment, err := c.cache.Get(ctx, id)
	if err != nil {
		return domain.Comment{}, err
	}
	if cacheComment.ID != 0 {
		return cacheComment, nil
	}
	// 再看数据库
	daoComment, err := c.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Comment{}, err
	}
	// 写入缓存
	err = c.cache.Set(ctx, domain.Comment(daoComment))
	if err != nil {
		log.Warning(err.Error())
	}
	return domain.Comment(daoComment), nil
}

func (c *commentRepo) GetByArticleID(ctx context.Context, id uint64, paging domain.Page, by domain.RangeBy) ([]domain.Comment, error) {
	// 查询缓存
	cacheComments, err := c.cache.ZGet(ctx, id, by)
	if err == nil {
		return cacheComments, nil
	}
	// 查询dao
	daoComments, err := c.dao.GetByArticleID(ctx, id, paging)
	if err != nil {
		return nil, err
	}
	comments := []domain.Comment{}
	// 组装数据
	for i, _ := range daoComments {
		comments = append(comments, domain.Comment(daoComments[i]))
	}
	// 写入缓存
	err = c.cache.ZAdd(ctx, comments...)
	if err != nil {
		log.Error(err.Error())
	}
	return comments, nil
}

func (c *commentRepo) GetByUserID(ctx context.Context, id uint64, paging domain.Page) ([]domain.Comment, error) {
	var comments []domain.Comment
	// 用户评论作为内部详细字段，暂无法从cache直接查询
	daoComments, err := c.dao.GetByUserID(ctx, id, paging)
	if err != nil {
		return nil, err
	}
	for i, _ := range daoComments {
		comments = append(comments, domain.Comment(daoComments[i]))
	}
	return comments, nil
}

func NewCommentRepo(dao dao.CommentDao, cache cache.CommentCache) CommentRepo {
	return &commentRepo{dao: dao, cache: cache}
}

func (c *commentRepo) dataBuild(comment model.Comment) domain.Comment {
	return domain.Comment{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		ArticleID: comment.ArticleID,
		ParentID:  comment.ParentID,
		Floor:     comment.Floor,
		State:     comment.State,
		Ctime:     comment.Ctime,
		Utime:     comment.Utime,
	}
}
