package cache

import (
	"douyin/package/constant"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func SetFavoriteSet(userID uint64, favoriteIDSet []uint64) error {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	favoriteIDStrings := make([]string, 1, len(favoriteIDSet)+1)
	favoriteIDStrings[0] = "0"
	for i := range favoriteIDSet {
		favoriteIDStrings = append(favoriteIDStrings, strconv.FormatUint(favoriteIDSet[i], 10))
	}
	pp := VideoRedisClient.Pipeline()
	pp.SAdd(key, favoriteIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

func GetFavoriteSet(userID uint64) ([]uint64, error) {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	// 若key不存在会返回空集合
	idSet, err := VideoRedisClient.SMembers(key).Result()
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
	err = VideoRedisClient.Del(videoInfoCountKey, videoInfoKey, favoriteKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
