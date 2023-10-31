package cache

import (
	"douyin/package/constant"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// TODO 为了简化代码 用zset的逻辑太复杂了
// 这里关注和点赞都用set，去redis拿到了set再去
// 数据库里查 也算用了缓存 也没有顺序问题
func SetFollowUserIDSet(userID uint64, followIDSet []uint64) error {
	key := constant.FollowIDPrefix + strconv.FormatUint(userID, 10)
	followIDStrings := make([]string, 0, len(followIDSet))
	for i := range followIDSet {
		followIDStrings = append(followIDStrings, strconv.FormatUint(followIDSet[i], 10))
	}
	pp := UserRedisClient.Pipeline()
	pp.SAdd(key, followIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

// 设置粉丝信息
func SetFollowerUserIDSet(userID uint64, followerIDSet []uint64) error {
	key := constant.FollowerIDPrefix + strconv.FormatUint(userID, 10)
	followIDStrings := make([]string, 0, len(followerIDSet))
	for i := range followerIDSet {
		followIDStrings = append(followIDStrings, strconv.FormatUint(followerIDSet[i], 10))
	}
	return UserRedisClient.SAdd(key, followIDStrings).Err()
}

func IsFollow(loginUserID, userID uint64) (bool, error) {
	key := constant.FollowIDPrefix + strconv.FormatUint(loginUserID, 10)
	// 应当判断键存不存再 不存在返回err
	exist, _ := UserRedisClient.Exists(key).Result()
	if exist > 0 {
		return UserRedisClient.SIsMember(key, strconv.FormatUint(userID, 10)).Result()
	}
	return false, redis.Nil
}

func GetFollowUserIDSet(userID uint64) ([]uint64, error) {
	key := constant.FollowIDPrefix + strconv.FormatUint(userID, 10)
	idSet, err := UserRedisClient.SMembers(key).Result()
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

func GetFollowerUserIDSet(userID uint64) ([]uint64, error) {
	key := constant.FollowerIDPrefix + strconv.FormatUint(userID, 10)
	idSet, err := UserRedisClient.SMembers(key).Result()
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

// lua保证原子性
// userID 点赞集合加1
// touserID 粉丝集合加1
// user表关注+1
// toUser 粉丝加1
func FollowAction(userID, toUserID uint64, cnt int64) error {
	followKey := constant.FollowIDPrefix + strconv.FormatUint(userID, 10)
	followerKey := constant.FollowerCountField + strconv.FormatUint(toUserID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	toUserInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(toUserID, 10)
	pp := UserRedisClient.Pipeline()
	if cnt == 1 {
		pp.SAdd(followKey, strconv.FormatUint(toUserID, 10))
		pp.SAdd(followerKey, strconv.FormatUint(userID, 10))
	} else {
		pp.SRem(followKey, strconv.FormatUint(toUserID, 10))
		pp.SRem(followerKey, strconv.FormatUint(userID, 10))
	}
	pp.HIncrBy(userInfoCountKey, constant.FollowCountField, cnt)
	pp.HIncrBy(toUserInfoCountKey, constant.FollowerCountField, cnt)
	_, err := pp.Exec()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

// // 使用Zset 排序使用数据库关注表的ID
// func SetFollowUserIDSet(userID uint64, followIDSet []Follow) error {
// 	pp := UserRedisClient.Pipeline()
// 	key := constant.FollowIDPrefix + strconv.FormatUint(userID, 10)
// 	zset := make([]redis.Z, len(followIDSet))
// 	for i := range followIDSet {
// 		zset[i] = redis.Z{
// 			Score:  float64(followIDSet[i].ID),
// 			Member: followIDSet[i].UserID,
// 		}
// 	}
// 	err1 := pp.ZAdd(key, zset...).Err()
// 	err2 := pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second).Err()
// 	_, err := pp.Exec()
// 	if err1 != nil {
// 		zap.L().Error(err1.Error())
// 	}
// 	if err2 != nil {
// 		zap.L().Error(err2.Error())
// 	}
// 	return err
// }
