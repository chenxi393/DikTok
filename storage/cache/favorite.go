package cache

import (
	"douyin/config"
	"douyin/package/constant"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var FavoriteRedisClient *redis.Client

func InitFavoriteRedis() {
	// FavoriteRedis 连接
	FavoriteRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.FavoriteRedis.Host, config.System.FavoriteRedis.Port),
		Password: config.System.FavoriteRedis.Password,
		DB:       config.System.FavoriteRedis.Database,
		PoolSize: config.System.FavoriteRedis.PoolSize, //每个CPU最大连接数
	})
	_, err := FavoriteRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("favorite_redis连接失败", zap.Error(err))
	}
}

// FIXME记得改造成 zset 需要按照点赞倒叙展示
func SetFavoriteSet(userID uint64, favoriteIDSet []uint64) error {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	favoriteIDStrings := make([]string, 1, len(favoriteIDSet)+1)
	favoriteIDStrings[0] = "0"
	for i := range favoriteIDSet {
		favoriteIDStrings = append(favoriteIDStrings, strconv.FormatUint(favoriteIDSet[i], 10))
	}
	pp := FavoriteRedisClient.Pipeline()
	pp.SAdd(key, favoriteIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

func GetFavoriteSet(userID uint64) ([]uint64, error) {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	// 若key不存在会返回空集合
	idSet, err := FavoriteRedisClient.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		return nil, redis.Nil
	}
	res := make([]uint64, 0, len(idSet))
	for _, t := range idSet {
		id, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

// 注意要更新redis 视频表的点赞数 点赞表 用户的点赞数 用户表的被点赞数（弃用）（目前采取删缓存）
func FavoriteAction(userID, author_id, videoID uint64, cnt int64) error {
	videoInfoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(videoID, 10)
	videoInfoKey := constant.VideoInfoPrefix + strconv.FormatUint(videoID, 10)
	favoriteKey := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(author_id, 10)
	err := UserRedisClient.Del(userInfoCountKey, authorInfoCountKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	err = VideoRedisClient.Del(videoInfoCountKey, videoInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	err = FavoriteRedisClient.Del(favoriteKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
