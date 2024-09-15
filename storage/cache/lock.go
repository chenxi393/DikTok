package cache

import (
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// 使用Redis的SetNX 实现分布式锁
// 尝试获取缓存前先加锁 加锁失败 等待10ms再去加锁
// 保证同一时间只有一个协程查数据库（不需要严格保证）
// 查完数据库后更新缓存然后释放锁
// 后续请求都会走缓存

const unlockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`

func GetLock(key, value string, exp time.Duration, client *redis.Client) (bool, error) {
	ok, err := client.SetNX(key, value, exp).Result()
	if err != nil {
		return false, err
	}
	return ok, err
}

func ReleaseLock(key, value string, client *redis.Client) error {
	script := redis.NewScript(unlockScript)
	err := script.Run(client, []string{key}, value).Err()
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}
