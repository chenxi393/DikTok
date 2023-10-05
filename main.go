package main

import (
	"douyin/config"
	"douyin/dal/dao"
	"douyin/package/cache"
	"douyin/package/util"
)

func main() {
	// 手动调用初始化函数
	config.Init()
	dao.InitMysql()
	cache.InitRedis()
	// TODO 消息队列
	util.InitZap()
	initFiber()

}
