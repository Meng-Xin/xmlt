package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"xmlt/global"
	"xmlt/internal/model"
	"xmlt/internal/shared"
	"xmlt/utils"
)

const (
	saveArticleState    = 0
	pendingArticleState = 1
	passArticleState    = 2
	deleteArticleState  = 3
)

type ArticleDAO interface {
	// Insert 的概念更加贴近关系型数据库，所以这里就不再是用 CREATE 这种说法了
	Insert(ctx context.Context, article model.Article) (uint64, error)
	Update(ctx context.Context, article model.Article) error
	UpdateToPublish(ctx context.Context, id uint64) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.Article, error)
	GetArticlesByCategoryID(ctx context.Context, categoryID uint64, paging *shared.Page) ([]model.Article, error)
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &articleGORM{db: db}
}

type articleGORM struct {
	db *gorm.DB
}

func (a *articleGORM) UpdateToPublish(ctx context.Context, id uint64) (uint64, error) {
	err := a.db.WithContext(ctx).Model(&model.Article{ID: id}).Update("state=?", pendingArticleState).Error
	return id, err
}

func (a *articleGORM) Insert(ctx context.Context, article model.Article) (uint64, error) {
	err := a.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

func (a *articleGORM) Update(ctx context.Context, article model.Article) error {
	article.Utime = utils.GetTimeMilli()
	return a.db.WithContext(ctx).Model(&article).Updates(article).Error
}

func (a *articleGORM) GetByID(ctx context.Context, id uint64) (model.Article, error) {
	var art model.Article
	err := a.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	return art, err
}

// AfterCreate 延时创建文章
func AfterCreate(db *gorm.DB) error {
	// 该端函数存在问题
	if db.Error != nil {
		return db.Error
	}
	//client, ok := db.Get("redis_client")
	client := global.Redis
	if client != nil {
		var article model.Article
		data, err := json.Marshal(&article)
		if err != nil {
			return err
		}
		result, err := client.Set(fmt.Sprintf("article_db_%d", article.ID), string(data), time.Hour).Result()
		if result != "OK" {
			return errors.New("插入失败")
		}
	}
	return nil
}

func (a *articleGORM) GetArticlesByCategoryID(ctx context.Context, categoryID uint64, page *shared.Page) ([]model.Article, error) {
	var articles []model.Article
	// 插入分页构造
	err := a.db.WithContext(ctx).Preload("User").Scopes(page.Paginate(&model.Article{})).Where("category_id=?", categoryID).Find(&articles).Error
	return articles, err
}
