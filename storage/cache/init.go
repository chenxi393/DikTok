package cache

import (
	"fmt"

	"diktok/config"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

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
