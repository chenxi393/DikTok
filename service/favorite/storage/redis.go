package storage

import (
	"math/rand"
	"strconv"
	"time"

	"diktok/package/constant"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	FavoriteRedis, VideoRedis, UserRedis *redis.Client
)

// FIXME记得改造成 zset 需要按照点赞倒叙展示
func SetFavoriteSet(userID int64, favoriteIDSet []int64) error {
	key := constant.FavoriteIDPrefix + strconv.FormatInt(userID, 10)
	favoriteIDStrings := make([]string, 1, len(favoriteIDSet)+1)
	favoriteIDStrings[0] = "0"
	for i := range favoriteIDSet {
		favoriteIDStrings = append(favoriteIDStrings, strconv.FormatInt(favoriteIDSet[i], 10))
	}
	pp := FavoriteRedis.Pipeline()
	pp.SAdd(key, favoriteIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

func GetFavoriteSet(userID int64) ([]int64, error) {
	key := constant.FavoriteIDPrefix + strconv.FormatInt(userID, 10)
	// 若key不存在会返回空集合
	idSet, err := FavoriteRedis.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		return nil, redis.Nil
	}
	res := make([]int64, 0, len(idSet))
	for _, t := range idSet {
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

// 注意要更新redis 视频表的点赞数 点赞表 用户的点赞数 用户表的被点赞数（弃用）（目前采取删缓存）
func FavoriteAction(userID, author_id, videoID int64, cnt int64) error {
	videoInfoCountKey := constant.VideoInfoCountPrefix + strconv.FormatInt(videoID, 10)
	videoInfoKey := constant.VideoInfoPrefix + strconv.FormatInt(videoID, 10)
	favoriteKey := constant.FavoriteIDPrefix + strconv.FormatInt(userID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatInt(userID, 10)
	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatInt(author_id, 10)
	err := UserRedis.Del(userInfoCountKey, authorInfoCountKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	err = VideoRedis.Del(videoInfoCountKey, videoInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	err = FavoriteRedis.Del(favoriteKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
