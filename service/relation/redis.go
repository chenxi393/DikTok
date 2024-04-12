package main

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

func SetFollowUserIDSet(userID int64, followIDSet []int64) error {
	key := constant.FollowIDPrefix + strconv.FormatInt(userID, 10)
	// 初始化的时候 加一个0 维持缓存存在
	followIDStrings := make([]string, 0, len(followIDSet))
	for i := range followIDSet {
		followIDStrings = append(followIDStrings, strconv.FormatInt(followIDSet[i], 10))
	}
	pp := relationRedis.Pipeline()
	pp.SAdd(key, followIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

// 设置粉丝信息
func SetFollowerUserIDSet(userID int64, followerIDSet []int64) error {
	key := constant.FollowerIDPrefix + strconv.FormatInt(userID, 10)
	followIDStrings := make([]string, 0, len(followerIDSet))
	for i := range followerIDSet {
		followIDStrings = append(followIDStrings, strconv.FormatInt(followerIDSet[i], 10))
	}
	pp := relationRedis.Pipeline()
	pp.SAdd(key, followIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

func IsFollow(loginUserID, userID int64) (bool, error) {
	key := constant.FollowIDPrefix + strconv.FormatInt(loginUserID, 10)
	// 应当判断键存不存再 不存在返回err
	exist, _ := relationRedis.Exists(key).Result()
	if exist > 0 {
		return relationRedis.SIsMember(key, strconv.FormatInt(userID, 10)).Result()
	}
	return false, redis.Nil
}

func GetFollowUserIDSet(userID int64) ([]int64, error) {
	key := constant.FollowIDPrefix + strconv.FormatInt(userID, 10)
	// 注意 若key不存在 则会返回空集合
	idSet, err := relationRedis.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		zap.L().Info(redis.Nil.Error())
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

func GetFollowerUserIDSet(userID int64) ([]int64, error) {
	key := constant.FollowerIDPrefix + strconv.FormatInt(userID, 10)
	idSet, err := relationRedis.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if len(idSet) == 0 {
		zap.L().Error(redis.Nil.Error())
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

// lua保证原子性 弃用 通通删缓存
// userID 点赞集合加1
// touserID 粉丝集合加1
// user表关注+1
// toUser 粉丝加1
func FollowAction(userID, toUserID int64, cnt int64) error {
	followKey := constant.FollowIDPrefix + strconv.FormatInt(userID, 10)
	followerKey := constant.FollowerIDPrefix + strconv.FormatInt(toUserID, 10)
	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatInt(userID, 10)
	toUserInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatInt(toUserID, 10)
	userInfoKey := constant.UserInfoPrefix + strconv.FormatInt(userID, 10)
	toUserInfoKey := constant.UserInfoPrefix + strconv.FormatInt(toUserID, 10)
	err := relationRedis.Del(followKey, followerKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	err = userRedis.Del(userInfoCountKey, toUserInfoCountKey, userInfoKey, toUserInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
