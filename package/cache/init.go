package cache

import (
	"douyin/config"
	"douyin/model"
	"douyin/package/constant"
	"fmt"
	"strconv"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var UserRedisClient *redis.Client
var VideoRedisClient *redis.Client
var CommentRedisClient *redis.Client

var UserIDBloomFilter *bloom.BloomFilter
var VideoIDBloomFilter *bloom.BloomFilter

func InitRedis() {
	// userRedis 连接
	UserRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.SystemConfig.UserRedis.Host, config.SystemConfig.UserRedis.Port),
		Password: config.SystemConfig.UserRedis.Password,
		DB:       config.SystemConfig.UserRedis.Database,
		PoolSize: config.SystemConfig.UserRedis.PoolSize, //每个CPU最大连接数
	})
	_, err := UserRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("user_redis连接失败", zap.Error(err))
	}
	// videoRedis 连接
	VideoRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.SystemConfig.VideoRedis.Host, config.SystemConfig.VideoRedis.Port),
		Password: config.SystemConfig.VideoRedis.Password,
		DB:       config.SystemConfig.VideoRedis.Database,
		PoolSize: config.SystemConfig.VideoRedis.PoolSize, //每个CPU最大连接数
	})
	_, err = VideoRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("video_redis连接失败", zap.Error(err))
	}
	// videoRedis 连接
	CommentRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.SystemConfig.CommentRedis.Host, config.SystemConfig.CommentRedis.Port),
		Password: config.SystemConfig.CommentRedis.Password,
		DB:       config.SystemConfig.CommentRedis.Database,
		PoolSize: config.SystemConfig.CommentRedis.PoolSize, //每个CPU最大连接数
	})
	_, err = VideoRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("comment_redis连接失败", zap.Error(err))
	}
	//
	zap.L().Info("redis连接成功成功")

	initBloomFilter()
}

// 初始化布隆过滤器
// 布隆过滤器的预估元素数量 和误报率 决定了底层bitmap的大小 和 无偏哈希函数的个数
func initBloomFilter() {
	UserIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	userIDList := make([]uint64, 0)
	constant.DB.Model(&model.User{}).Select("id").Find(&userIDList)
	for _, u := range userIDList {
		UserIDBloomFilter.AddString(strconv.FormatUint(u, 10))
	}
	VideoIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	videoIDList := make([]uint64, 0)
	constant.DB.Model(&model.Video{}).Select("id").Find(&videoIDList)
	for _, v := range videoIDList {
		VideoIDBloomFilter.AddString(strconv.FormatUint(v, 10))
	}
	zap.L().Info("初始化布隆过滤器成功")
}
