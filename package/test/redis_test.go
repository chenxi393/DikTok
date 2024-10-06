package test

import (
	"diktok/config"
	"diktok/storage/cache"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestMget(t *testing.T) {
	viper.AddConfigPath("../../config/")
	config.Init()
	rdb := cache.InitRedis(config.System.Redis.CommentDB)
	// 写入数据
	err := rdb.Set("key1", "value1", 0).Err()
	if err != nil {
		fmt.Println("Error setting key1:", err)
	}

	err = rdb.Set("key2", "value2", 0).Err()
	if err != nil {
		fmt.Println("Error setting key2:", err)
	}

	err = rdb.Set("key3", "value3", 0).Err()
	if err != nil {
		fmt.Println("Error setting key3:", err)
	}

	// 使用 MGet 查询
	values, err := rdb.MGet("dddd", "key1", "sddd", "key2", "key3", "dddd").Result()
	if err != nil {
		fmt.Println("Error getting keys:", err)
	} else {
		fmt.Println("Values:", values)
	}
}
