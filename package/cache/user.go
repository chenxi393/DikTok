package cache

import (
	"douyin/config"
	"douyin/model"
	"douyin/package/constant"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var UserRedisClient *redis.Client
var UserIDBloomFilter *bloom.BloomFilter

// redis的用户表
// 用户信息拆成两部分
// 分为一般不轻易改变的字段UserInfo
// 容易改变的字段UserInfoCount
// 目的: 若直接存储在hash里 会频繁变动性能问题

// UserInfo 用户信息中基本信息，不做更改或更改频率较低
type UserInfo struct {
	ID              uint64 // 自增主键
	Username        string // 用户名
	Password        string // 用户密码
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	Signature       string // 个人简介
}

func InitUserRedis() {
	// userRedis 连接
	UserRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.UserRedis.Host, config.System.UserRedis.Port),
		Password: config.System.UserRedis.Password,
		DB:       config.System.UserRedis.Database,
		PoolSize: config.System.UserRedis.PoolSize, //每个CPU最大连接数
	})
	_, err := UserRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("user_redis连接失败", zap.Error(err))
	}
	initUserBloomFilter()
}

func initUserBloomFilter() {
	// 估计会有10万个用户 误报率是0.01
	UserIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	userIDList := make([]uint64, 0)
	constant.DB.Model(&model.User{}).Select("id").Find(&userIDList)
	for _, u := range userIDList {
		UserIDBloomFilter.AddString(strconv.FormatUint(u, 10))
	}
}

// 拆分userInfo 哈希存储易变的字段 string存储基本不变的字段
// 使用pipeline 一次执行
func SetUserInfo(user *model.User) error {
	// 拆分成两个数据结构，更改频率较低放在 UserInfo 中
	userInfo := &UserInfo{
		ID:              user.ID,
		Username:        user.Username,
		Password:        user.Password,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
	}
	// 进行序列化
	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	// 开启管道 一次发送请求
	pipeline := UserRedisClient.Pipeline()

	// 下面两个的过期时间保持一致 不然查库还是会查出信息
	randomTime := rand.Intn(100)
	// 设置 UserInfo 的 JSON 缓存
	infoKey := constant.UserInfoPrefix + strconv.FormatUint(user.ID, 10)
	err = pipeline.Set(infoKey, userInfoJSON,
		constant.Expiration+time.Duration(randomTime)*time.Second).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}

	infoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(user.ID, 10)
	// 使用 MSet 进行批量设置
	err = pipeline.HMSet(infoCountKey, map[string]interface{}{
		constant.FollowCountField:    user.FollowCount,
		constant.FollowerCountField:  user.FollowerCount,
		constant.TotalFavoritedField: user.TotalFavorited,
		constant.WorkCountField:      user.WorkCount,
		constant.FavoriteCountField:  user.FavoriteCount,
	}).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	err = pipeline.Expire(infoCountKey, constant.Expiration+time.Duration(randomTime)*time.Second).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	// 执行管道中的命令
	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	return nil
}

func GetUserInfo(userID uint64) (*model.User, error) {
	infoKey := constant.UserInfoPrefix + strconv.FormatUint(userID, 10)
	infoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	// 使用管道加速
	pipeline := UserRedisClient.Pipeline()
	// 注意pipeline返回指针 返回值肯定是nil
	userInfoCmd := pipeline.Get(infoKey)               // key 不存在会返回 redis:nil err
	userInfoCountCmd := pipeline.HGetAll(infoCountKey) // 注意这里键不存在又不会返回nil err
	_, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	// 提取返回的结果
	userInfo, err := userInfoCmd.Result()
	if err != nil {
		return nil, err
	}
	userInfoCount, err := userInfoCountCmd.Result()
	if err != nil {
		return nil, err
	}
	if len(userInfoCount) == 0 {
		err = errors.New("")
		return nil, err
	}
	// 解析不变的字段
	userInfoFixed := UserInfo{}
	err = json.Unmarshal([]byte(userInfo), &userInfoFixed)
	if err != nil {
		return nil, err
	}

	// 解析count信息
	followCount, _ := strconv.ParseInt(userInfoCount[constant.FollowCountField], 10, 64)
	followerCount, _ := strconv.ParseInt(userInfoCount[constant.FollowerCountField], 10, 64)
	totalFavoritedCount, _ := strconv.ParseInt(userInfoCount[constant.TotalFavoritedField], 10, 64)
	workCount, _ := strconv.ParseInt(userInfoCount[constant.WorkCountField], 10, 64)
	favoriteCount, _ := strconv.ParseInt(userInfoCount[constant.FavoriteCountField], 10, 64)

	return &model.User{
		ID:              userInfoFixed.ID,
		Username:        userInfoFixed.Username,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		Avatar:          userInfoFixed.Avatar,
		BackgroundImage: userInfoFixed.BackgroundImage,
		Signature:       userInfoFixed.Signature,
		TotalFavorited:  totalFavoritedCount,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}, nil
}
