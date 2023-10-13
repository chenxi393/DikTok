package main

import (
	"douyin/config"
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/util"
	"douyin/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 手动调用初始化函数
	app := fiber.New(fiber.Config{
		BodyLimit: 40 * 1024 * 1024, // 不设置这个会返回413 客户端文件太大 导致返回413
	})

	config.Init()
	database.InitMysql()
	cache.InitRedis()
	// TODO 消息队列
	util.InitZap()

	app.Use(logger.New()) // 使用中间件打印日志
	router.InitRouter(app)
	panic(app.Listen(
		config.SystemConfig.HttpAddress.Host + ":" + config.SystemConfig.HttpAddress.Port).Error())
}
