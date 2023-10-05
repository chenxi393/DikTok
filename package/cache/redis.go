package cache

import (
	"douyin/config"

	"github.com/go-redis/redis"
)

var UserRedisClient *redis.Client

func InitRedis() {
	cliet := redis.NewClient(&redis.Options{
		Addr:     config.SystemConfig.Redis.Host + ":" + config.SystemConfig.Redis.Port,
		Password: config.SystemConfig.Redis.Password,
		DB:       config.SystemConfig.Redis.Database,
		PoolSize: config.SystemConfig.Redis.PoolSize, //每个CPU最大连接数
	})
	_, err := cliet.Ping().Result()
	if err != nil {
		panic(err.Error())
	}
	UserRedisClient = cliet
}
