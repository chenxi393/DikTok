package main

import (
	"douyin/config"
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/llm"
	"douyin/package/mq"
	"douyin/package/util"
	"douyin/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func main() {
	// 手动调用初始化函数 可以考虑使用init函数
	config.Init()
	util.InitZap()
	database.InitMySQL()
	cache.InitRedis()
	mq.InitMQ()
	llm.RegisterChatGPT()
	// 客户端文件超过30MB 返回413
	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})
	// 使用中间件打印日志
	app.Use(logger.New())
	router.InitRouter(app)
	zap.L().Fatal("fiber启动失败: ", zap.Error(app.Listen(
		config.System.HttpAddress.Host+":"+config.System.HttpAddress.Port)))
}
