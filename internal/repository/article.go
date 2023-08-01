package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
	"xmlt/internal/domain"
	"xmlt/internal/model"
	"xmlt/internal/repository/cache"
	"xmlt/internal/repository/dao"
	"xmlt/internal/shared"
	"xmlt/utils"
)

type ArticleRepo interface {
	// Create 创建一篇文章
	Create(ctx context.Context, article domain.Article) (uint64, error)
	CreateAndCached(ctx context.Context, article domain.Article) (uint64, error)
	// Update 更新一篇文章
	Update(ctx context.Context, article domain.Article) error
	// Get 方法应该负责将 Author 也一并组装起来。
	// 这里就会有一个很重要的概念，叫做延迟加载，但是 GO 是做不了的，
	// 所以只能考虑传递标记位，或者使用新方法来控制要不要把 Author 组装起来
	Get(ctx context.Context, id uint64) (domain.Article, error)
	GetArticlesByCategoryID(ctx context.Context, categoryID uint64, paging *shared.Page) ([]domain.Article, error)
}

func NewArticleRepo(dao dao.ArticleDAO, artCache cache.ArticleCache) ArticleRepo {
	return &articleRepo{dao: dao, cache: artCache}
}

type articleRepo struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func (a *articleRepo) Create(ctx context.Context, article domain.Article) (uint64, error) {
	now := utils.GetTimeMilli()
	entity := model.Article{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		UserID:       article.Author,
		CategoryID:   article.CategoryID,
		NiceTopic:    article.NiceTopic,
		BrowseCount:  article.BrowseCount,
		ThumbsUP:     article.ThumbsUP,
		Ctime:        now,
		Utime:        now,
	}
	return a.dao.Insert(ctx, entity)
}

func (a *articleRepo) CreateAndCached(ctx context.Context, article domain.Article) (uint64, error) {
	now := utils.GetTimeMilli()
	entity := model.Article{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		UserID:       article.Author,
		CategoryID:   article.CategoryID,
		NiceTopic:    article.NiceTopic,
		BrowseCount:  article.BrowseCount,
		ThumbsUP:     article.ThumbsUP,
		Ctime:        now,
		Utime:        now,
	}
	// 也可以再次封装error
	id, err := a.dao.Insert(ctx, entity)
	if err != nil {
		return 0, err
	}
	article.ID = id
	err = a.cache.Set(ctx, article)
	return id, err
}

func (a *articleRepo) Update(ctx context.Context, article domain.Article) error {
	now := utils.GetTimeMilli()
	entity := model.Article{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		UserID:       article.Author,
		CategoryID:   article.CategoryID,
		NiceTopic:    article.NiceTopic,
		BrowseCount:  article.BrowseCount,
		ThumbsUP:     article.ThumbsUP,
		Utime:        now,
	}
	return a.dao.Update(ctx, entity)
}

func (a *articleRepo) Get(ctx context.Context, id uint64) (domain.Article, error) {
	// 先看缓存，缓存命中直接返回
	res, err := a.cache.Get(ctx, id)
	if err == nil {
		return res, err
	}
	// 缓存未命中，执行数据库查询，并重新更新Redis
	entity, err := a.dao.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 按道理来说这里需要提前组装好 Author的
	art := domain.Article{
		ID:           entity.ID,
		Title:        entity.Title,
		Content:      entity.Content,
		CommentCount: entity.CommentCount,
		Status:       entity.Status,
		Author:       entity.UserID,
		CategoryID:   entity.CategoryID,
		NiceTopic:    entity.NiceTopic,
		BrowseCount:  entity.BrowseCount,
		ThumbsUP:     entity.ThumbsUP,
		Ctime:        time.UnixMilli(entity.Ctime),
		Utime:        time.UnixMilli(entity.Utime),
	}
	err = a.cache.Set(ctx, art)
	if err != nil {
		// 这个 error 可以考虑是否忽略，一位内缓存虽然更新失败了，但实际上数据库已经拿到了
		// 不过实际上这个很危险，因为如果 Redis 整个崩溃了,那么数据库肯定是扛不住压力的
		log.Error(err)
	}
	return art, nil
}

func (a *articleRepo) GetArticlesByCategoryID(ctx context.Context, categoryID uint64, paging *shared.Page) ([]domain.Article, error) {
	// 缓存未命中，执行数据库查询，并重新更新Redis
	data, err := a.dao.GetArticlesByCategoryID(ctx, categoryID, paging)
	if err != nil {
		return nil, err
	}
	var entities []domain.Article
	for i, _ := range data {
		// 按道理来说这里需要提前组装好 Author的
		entity := domain.Article{
			ID:          data[i].ID,
			Title:       data[i].Title,
			NiceTopic:   data[i].NiceTopic,
			BrowseCount: data[i].BrowseCount,

			Ctime: time.UnixMilli(data[i].Ctime),
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func (a *articleRepo) dao2Dto(entity model.Article) domain.Article {
	domainArticle := domain.Article{
		ID:           entity.ID,
		Title:        entity.Title,
		Content:      entity.Content,
		CommentCount: entity.CommentCount,
		Status:       entity.Status,
		CategoryID:   entity.CategoryID,
		NiceTopic:    entity.NiceTopic,
		BrowseCount:  entity.BrowseCount,
		ThumbsUP:     entity.ThumbsUP,
		Ctime:        time.UnixMilli(entity.Ctime),
		Utime:        time.UnixMilli(entity.Utime),
	}
	// 将DAO对象的切片字段转换为Domain对象的切片字段
	for _, commentDAO := range entity.Comments {
		domainComment := domain.Comment{
			ID:       commentDAO.ID,
			Content:  commentDAO.Content,
			UserID:   commentDAO.UserID,
			ParentID: commentDAO.ParentID,
			Floor:    commentDAO.Floor,
		}
		domainArticle.Comments = append(domainArticle.Comments, domainComment)
	}
	for _, likeDao := range entity.UserLikes {
		domainUserLike := domain.UserLikeArticle{
			ID:     likeDao.ID,
			UserID: likeDao.UserID,
		}
		domainArticle.UserLikes = append(domainArticle.UserLikes, domainUserLike)
	}
	return domainArticle
}
