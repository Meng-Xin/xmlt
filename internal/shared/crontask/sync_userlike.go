package crontask

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"xmlt/global"
	"xmlt/internal/model"
	"xmlt/internal/shared/enum"
)

type UserLikeCron interface {
	SyncFromRedisToMysql()
}

func NewUserLikeCron(scheduler *gocron.Scheduler, client *redis.Client, db *gorm.DB) UserLikeCron {
	return &userLikeCron{schedule: scheduler, client: client, db: db}
}

type userLikeCron struct {
	schedule *gocron.Scheduler
	client   *redis.Client
	db       *gorm.DB
}

func (u *userLikeCron) SyncFromRedisToMysql() {
	// 每 5 个小时从Redis同步到Mysql
	u.schedule.Every(5).Hour().Do(func() {
		ctx := context.Background()
		keys, err := u.client.WithContext(ctx).Keys("*").Result()
		if err != nil {
			panic(err)
		}
		// 匹配符合要求的key
		matchKey := make([]string, 0)
		for i, _ := range keys {
			if ok := strings.Contains(keys[i], enum.UserLikeArticle[:18]); ok {
				matchKey = append(matchKey, keys[i])
			}
		}

		// 遍历Redis集合并更新到MySQL
		for i, _ := range matchKey {
			endKey := strings.Split(matchKey[i], "_")
			articleId, _ := strconv.ParseUint(endKey[3], 10, 64)
			members, err := u.client.WithContext(ctx).ZRangeWithScores(matchKey[i], 0, -1).Result()

			if err != nil {
				panic(err)
			}
			for _, member := range members {
				userId, _ := strconv.ParseUint(member.Member.(string), 10, 64)
				score := member.Score
				fmt.Printf("Member: %s, Score: %f\n", userId, score)
				// 更新到Mysql
				//userId, _ := strconv.ParseUint(memberValue, 10, 64)
				entity := model.UserLikeArticle{
					ArticleID: articleId,
					UserID:    userId,
					LikeState: true,
					Ctime:     uint64(score),
					Utime:     uint64(score),
				}

				// 开启事务从Redis更新数据到Mysql
				u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
					daoUser := model.UserLikeArticle{}
					err = tx.WithContext(ctx).Where("article_id=? and user_id=?", articleId, userId).First(&daoUser).Error
					if err == nil && daoUser.ID != 0 {
						//  TODO ArticleId相同前提下的所有 UserID，与Redis存储的Id取差集，把MySQL的点赞状态给取消掉；
						err = tx.WithContext(ctx).Where("article_id=? and user_id=?", articleId, userId).Updates(&entity).Error
						if err != nil {
							panic(err)
						}
					} else {
						// 执行创建
						err = tx.WithContext(ctx).Model(model.UserLikeArticle{}).Where("article_id=? and user_id=?", articleId, userId).Save(&entity).Error
						if err != nil {
							panic(err)
						}
					}
					return nil
				})
			}

		}
	})

}

type user struct {
	id    uint64 // 用户id
	score uint64 // 记录时间
}

// GetRedisMembers 从Redis获取点赞文章对应的结构
func (u *userLikeCron) GetRedisMembers(ctx context.Context, client *redis.Client) map[uint64][]user {
	redisMark := make(map[uint64][]user, 0) // key: ArticleID Val: userIDs

	keys, err := u.client.WithContext(ctx).Keys("*").Result()
	if err != nil {
		panic(err)
	}
	// 匹配符合要求的key user_like_article_
	matchKey := make([]string, 0)
	for i, _ := range keys {
		if ok := strings.Contains(keys[i], enum.UserLikeArticle[:18]); ok {
			matchKey = append(matchKey, keys[i])
		}
	}

	// 遍历Redis集合并更新到MySQL %d
	for i, _ := range matchKey {
		endKey := strings.Split(matchKey[i], "_")
		articleId, _ := strconv.ParseUint(endKey[3], 10, 64)
		members, err := u.client.WithContext(ctx).ZRangeWithScores(matchKey[i], 0, -1).Result()
		if err != nil {
			panic(err)
		}
		users := make([]user, 0)
		for _, member := range members {
			userId, _ := strconv.ParseUint(member.Member.(string), 10, 64)
			users = append(users, user{id: userId, score: uint64(member.Score)})
		}
		redisMark[articleId] = users
	}
	return redisMark
}

// GetMysqlMembers 从Mysql获取点赞用户信息
func (u *userLikeCron) GetMysqlMembers(ctx context.Context, artId uint64) []uint64 {
	mysqlMark := make([]uint64, 0) // key: ArticleID Val: userIDs
	daoUserLikes := []model.UserLikeArticle{}
	err := u.db.WithContext(ctx).Where("article_id=?", artId).Find(&daoUserLikes).Error
	if err != nil {
		global.Log.Warn(err.Error())
		return nil
	}
	for _, daoUser := range daoUserLikes {
		mysqlMark = append(mysqlMark, daoUser.UserID)
	}
	return mysqlMark
}
