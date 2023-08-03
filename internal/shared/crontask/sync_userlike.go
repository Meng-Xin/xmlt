package crontask

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"xmlt/global"
	"xmlt/internal/model"
	"xmlt/internal/shared/enum"
	"xmlt/utils"
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

// SyncFromRedisToMysql 触发条件：1.Redis内至少存在一条有关点赞的记录信息，否则不加载。
func (u *userLikeCron) SyncFromRedisToMysql() {
	// 每 5 个小时从Redis同步到Mysql
	u.schedule.Every(30).Second().Do(func() {
		ctx := context.Background()
		redisMarks, err := u.GetRedisMembers(ctx)
		if err != nil {
			global.Log.Warn(err.Error())
		}
		for artId, val := range redisMarks {
			// key：uid Val: time ,本列表中的点赞信息记录
			likeInfo := make(map[uint64]uint64, 0)
			redisUsers := []uint64{}
			mysqlUsers, err := u.GetMysqlMembers(ctx, artId)
			if err != nil {
				global.Log.Warn(err.Error())
			}
			for i, _ := range val {
				// 获取完整uid列表
				redisUsers = append(redisUsers, val[i].id)
				// 拼接时间信息
				likeInfo[val[i].id] = val[i].score
			}
			// 已存在于Mysql但是不存在于Redis
			comparedRefMysql := utils.DifferenceCompared(mysqlUsers, redisUsers)
			// 已存在于Redis但是不存在于MySQL
			comparedRefRedis := utils.DifferenceCompared(redisUsers, mysqlUsers)
			// comparedRefMysql 差集取消用户点赞状态
			for _, uid := range comparedRefMysql {
				err = u.UpdateUserLikeState(ctx, artId, uid, likeInfo[uid])
				if err != nil {
					log.Warn(err.Error())
				}
			}
			// comparedRefRedis 差集补充用户点赞信息
			for _, uid := range comparedRefRedis {
				err = u.CreateUserLikeState(ctx, artId, uid, likeInfo[uid])
				if err != nil {
					log.Warn(err.Error())
				}
			}
		}
	})

}

type user struct {
	id    uint64 // 用户id
	score uint64 // 记录时间
}

// GetRedisMembers 从Redis获取点赞文章对应的结构
func (u *userLikeCron) GetRedisMembers(ctx context.Context) (map[uint64][]user, error) {
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
			return nil, err
		}
		users := make([]user, 0)
		for _, member := range members {
			userId, _ := strconv.ParseUint(member.Member.(string), 10, 64)
			users = append(users, user{id: userId, score: uint64(member.Score)})
		}
		redisMark[articleId] = users
	}
	return redisMark, nil
}

// GetMysqlMembers 从Mysql获取点赞用户信息
func (u *userLikeCron) GetMysqlMembers(ctx context.Context, artId uint64) ([]uint64, error) {
	mysqlMark := make([]uint64, 0) // key: ArticleID Val: userIDs
	daoUserLikes := []model.UserLikeArticle{}
	err := u.db.WithContext(ctx).Where("article_id=?", artId).Find(&daoUserLikes).Error
	if err != nil {
		return nil, err
	}
	for _, daoUser := range daoUserLikes {
		mysqlMark = append(mysqlMark, daoUser.UserID)
	}
	return mysqlMark, nil
}

func (u *userLikeCron) UpdateUserLikeState(ctx context.Context, artId, uid uint64, timer uint64) error {
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where("article_id=? and user_id=?", artId, uid).Updates(&model.UserLikeArticle{LikeState: false, Utime: timer}).Error
		return err
	})

	return err
}

func (u *userLikeCron) CreateUserLikeState(ctx context.Context, artId, uid uint64, timer uint64) error {
	entity := model.UserLikeArticle{
		ArticleID: artId,
		UserID:    uid,
		LikeState: true,
		Ctime:     timer,
	}
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(model.UserLikeArticle{}).Create(&entity).Error
		return err
	})
	return err
}
