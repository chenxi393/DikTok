package cache

import (
	"douyin/config"
	"douyin/model"
	"douyin/package/database"
	"fmt"
	"strconv"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var UserIDBloomFilter *bloom.BloomFilter
var VideoIDBloomFilter *bloom.BloomFilter

func InitRedis(db int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.Redis.Host, config.System.Redis.Port),
		Password: config.System.Redis.Password,
		DB:       db,
		PoolSize: config.System.Redis.PoolSize, //每个CPU最大连接数
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("redis连接失败", zap.Error(err))
	}
	if db == config.System.Redis.UserDB {
		initUserBloomFilter()
	} else if db == config.System.Redis.VideoDB {
		initVideoBloomFilter()
	}
	return redisClient
}

func initUserBloomFilter() {
	if database.DB == nil {
		database.InitMySQL()
	}
	// 估计会有10万个用户 误报率是0.01
	UserIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	userIDList := make([]uint64, 0)
	database.DB.Model(&model.User{}).Select("id").Find(&userIDList)
	for _, u := range userIDList {
		UserIDBloomFilter.AddString(strconv.FormatUint(u, 10))
	}
}

// 初始化布隆过滤器
// 布隆过滤器的预估元素数量 和误报率 决定了底层bitmap的大小 和 无偏哈希函数的个数
func initVideoBloomFilter() {
	if database.DB == nil {
		database.InitMySQL()
	}
	VideoIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	videoIDList := make([]uint64, 0)
	database.DB.Model(&model.Video{}).Select("id").Find(&videoIDList)
	for _, v := range videoIDList {
		VideoIDBloomFilter.AddString(strconv.FormatUint(v, 10))
	}
	zap.L().Info("初始化布隆过滤器: 成功")
}
