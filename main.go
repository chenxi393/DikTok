package main

import (
	"douyin/config"
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/llm"
	"douyin/package/mq"
	"douyin/package/util"
)

func main() {
	// TODO 配置文件实际上也应该分离
	config.Init()
	// TODO 日志也应该考虑合并
	util.InitZap()
	database.InitMySQL()
	cache.InitRedis()
	mq.InitMQ()
	llm.RegisterChatGPT()
}
