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
		Addr:     fmt.Sprintf("%s:%s", config.System.UserRedis.Host, config.System.UserRedis.Port),
		Password: config.System.UserRedis.Password,
		DB:       config.System.UserRedis.Database,
		PoolSize: config.System.UserRedis.PoolSize, //每个CPU最大连接数
	})
	_, err := UserRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("user_redis连接失败", zap.Error(err))
	}
	// videoRedis 连接
	VideoRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.VideoRedis.Host, config.System.VideoRedis.Port),
		Password: config.System.VideoRedis.Password,
		DB:       config.System.VideoRedis.Database,
		PoolSize: config.System.VideoRedis.PoolSize, //每个CPU最大连接数
	})
	_, err = VideoRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("video_redis连接失败", zap.Error(err))
	}
	// videoRedis 连接
	CommentRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.CommentRedis.Host, config.System.CommentRedis.Port),
		Password: config.System.CommentRedis.Password,
		DB:       config.System.CommentRedis.Database,
		PoolSize: config.System.CommentRedis.PoolSize, //每个CPU最大连接数
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
