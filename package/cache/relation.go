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
	// 初始化的时候 加一个0 维持缓存存在
	followIDStrings := make([]string, 1, len(followIDSet)+1)
	followIDStrings[0] = "0"
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
	followIDStrings := make([]string, 1, len(followerIDSet)+1)
	followIDStrings[0] = "0"
	for i := range followerIDSet {
		followIDStrings = append(followIDStrings, strconv.FormatUint(followerIDSet[i], 10))
	}
	pp := UserRedisClient.Pipeline()
	pp.SAdd(key, followIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
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
	// 注意 若key不存在 则会返回空集合
	idSet, err := UserRedisClient.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		zap.L().Info(redis.Nil.Error())
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

func GetFollowerUserIDSet(userID uint64) ([]uint64, error) {
	key := constant.FollowerIDPrefix + strconv.FormatUint(userID, 10)
	idSet, err := UserRedisClient.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		zap.L().Error(redis.Nil.Error())
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

// lua保证原子性 弃用 通通删缓存
// userID 点赞集合加1
// touserID 粉丝集合加1
// user表关注+1
// toUser 粉丝加1
func FollowAction(userID, toUserID uint64, cnt int64) error {
	followKey := constant.FollowIDPrefix + strconv.FormatUint(userID, 10)
	followerKey := constant.FollowerIDPrefix + strconv.FormatUint(toUserID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	toUserInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(toUserID, 10)
	userInfoKey := constant.UserInfoPrefix + strconv.FormatUint(userID, 10)
	toUserInfoKey := constant.UserInfoPrefix + strconv.FormatUint(toUserID, 10)
	err := UserRedisClient.Del(followKey, followerKey, userInfoCountKey, toUserInfoCountKey, userInfoKey, toUserInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
