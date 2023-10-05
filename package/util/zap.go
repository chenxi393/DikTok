package util

import (
	"douyin/config"

	"go.uber.org/zap"
)

func InitZap() {
	// 先打印到控制台吧
	var logger *zap.Logger
	if config.SystemConfig.Mode == "debug" {
		logger, _ = zap.NewDevelopment()
	} else if config.SystemConfig.Mode == "example" {
		logger = zap.NewExample()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()
	zap.ReplaceGlobals(logger) //返回值似乎是一个取消函数
}
