package cache

import (
	"douyin/package/constant"
	"strconv"

	"go.uber.org/zap"
)

func SetFavoriteSet(userID uint64, favoriteIDSet []uint64) error {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	favoriteIDStrings := make([]string, 0, len(favoriteIDSet))
	for i := range favoriteIDSet {
		favoriteIDStrings = append(favoriteIDStrings, strconv.FormatUint(favoriteIDSet[i], 10))
	}
	return VideoRedisClient.SAdd(key, favoriteIDStrings).Err()
}

func GetFavoriteSet(userID uint64) ([]uint64, error) {
	key := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	idSet, err := VideoRedisClient.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
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

// FIXME 注意要更新redis 视频表的点赞数 点赞表 用户的点赞数 用户表的被点赞数
func FavoriteAction(userID, author_id, videoID uint64, cnt int64) error {
	videoInfoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(videoID, 10)
	favoriteKey := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(author_id, 10)
	pp := UserRedisClient.Pipeline()
	if cnt == 1 {
		pp.SAdd(favoriteKey, strconv.FormatUint(videoID, 10))
	} else {
		pp.SRem(favoriteKey, strconv.FormatUint(videoID, 10))
	}
	pp.HIncrBy(videoInfoCountKey, constant.FavoritedCountField, cnt)
	pp.HIncrBy(userInfoCountKey, constant.FavoriteCountField, cnt)
	pp.HIncrBy(authorInfoCountKey, constant.TotalFavoritedField, cnt)
	_, err := pp.Exec()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

// // 注意使用了Zset 因为关注是有顺序的 最近关注的在上 顺序按照数据库的ID排序
// func SetFavoriteVideoIDSet(userID uint64, favorites []Favorite) error {
// 	pp := VideoRedisClient.Pipeline()
// 	favoriteKey := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
// 	zset := make([]redis.Z, len(favorites))
// 	for i := range favorites {
// 		zset[i] = redis.Z{
// 			Score:  float64(favorites[i].ID),
// 			Member: favorites[i].VideoID,
// 		}
// 	}
// 	err1 := pp.ZAdd(favoriteKey, zset...).Err()
// 	err2 := pp.Expire(favoriteKey, constant.Expiration+time.Duration(rand.Intn(100))*time.Second).Err()
// 	_, err := pp.Exec()
// 	if err1 != nil {
// 		zap.L().Error(err1.Error())
// 	}
// 	if err2 != nil {
// 		zap.L().Error(err2.Error())
// 	}
// 	return err
// }
